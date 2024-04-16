/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const BaseUrl string = "https://opendev.org/openstack/openstack-helm/raw/branch/master"

var Cwd string
var Download bool
var OverridesPath string
var DownloadBaseUrl string

var rootCmd = &cobra.Command{
	Use:  "get-values-overrides",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chart := args[0]
		features := args[1:]
		overrideCandidates := generateOverrideCandidates(chart, features)
		overrideArgs := getOverrideHelmArguments(chart, overrideCandidates)
		fmt.Println(strings.Join(overrideArgs, " "))
	},
}

func init() {
	Cwd, _ = os.Getwd()
	rootCmd.Flags().BoolVarP(&Download, "download", "d", false, "Download the overrides from the internet if does not exist in the path (default: false)")
	rootCmd.Flags().StringVarP(&DownloadBaseUrl, "url", "u", BaseUrl, fmt.Sprintf("Base url to download overrides (default: %s)", BaseUrl))
	rootCmd.Flags().StringVarP(&OverridesPath, "path", "p", Cwd, "Path to the overrides (default: current directory)")
}

func num2items(num uint32, power int) []int {
	featureNums := make([]int, 0)
	for i := 0; i < power; i++ {
		if uint32(math.Pow(2, float64(i)))&num != 0 {
			featureNums = append(featureNums, i)
		}
	}
	slices.Reverse(featureNums)
	return featureNums
}

func overrideFile(chart, overrideName string) string {
	return filepath.Join(OverridesPath, chart, "values_overrides", overrideName)
}

func downloadOverride(chart, overrideName string) error {
	fullUrl := fmt.Sprintf("%s/%s/values_overrides/%s", DownloadBaseUrl, chart, overrideName)
	filename := overrideFile(chart, overrideName)
	fmt.Fprintf(os.Stderr, "Trying to download %s\n", fullUrl)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download %s: %s", fullUrl, resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func generateOverrideCandidates(chart string, features []string) []string {
	fmt.Fprintf(os.Stderr, "Chart: %s Features: %s\n", chart, strings.Join(features, " "))
	slices.Reverse(features)
	overrideCandidates := make([]string, 0)
	for num := uint32(1); num < uint32(math.Pow(2, float64(len(features)))); num++ {
		words := make([]string, 0)
		items := num2items(num, len(features))
		for _, item := range items {
			words = append(words, features[item])
		}
		overrideCandidates = append(overrideCandidates, fmt.Sprintf("%s.yaml", strings.Join(words, "-")))
	}
	return overrideCandidates
}

func getOverrideHelmArguments(chart string, overrideCandidates []string) []string {
	overrides := make([]string, 0)
	for _, overrideCandidate := range overrideCandidates {
		overrideCandidateFile := overrideFile(chart, overrideCandidate)
		fmt.Fprintf(os.Stderr, "Override candidate: %s\n", overrideCandidateFile)
		_, err := os.Stat(overrideCandidateFile)
		if err == nil {
			fmt.Fprintf(os.Stderr, "File found: %s\n", overrideCandidateFile)
			overrides = append(overrides, fmt.Sprintf("--values %s", overrideCandidateFile))
		} else {
			if Download {
				fmt.Fprintf(os.Stderr, "File not found: %s\n", overrideCandidateFile)
				err = downloadOverride(chart, overrideCandidate)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Fprintf(os.Stderr, "Successfully downloaded %s\n", overrideCandidate)
					overrides = append(overrides, fmt.Sprintf("--values %s", overrideCandidateFile))
				}
			}
		}
	}
	fmt.Fprintln(os.Stderr, "Resulting override Helm arguments:")
	fmt.Fprintln(os.Stderr, strings.Join(overrides, " "))
	return overrides
}

func main() {
	rootCmd.Execute()
}
