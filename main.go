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
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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

func main() {
	overridesPath := os.Args[1]
	fmt.Fprintln(os.Stderr, "Lookup path:", overridesPath)
	features := make([]string, len(os.Args[2:]))
	copy(features, os.Args[2:])
	fmt.Fprintf(os.Stderr, "Features: %s\n", strings.Join(features, " "))
	slices.Reverse(features)
	overrideCandidates := make([]string, 0)
	for num := uint32(1); num < uint32(math.Pow(2, float64(len(features)))); num++ {
		words := make([]string, 0)
		items := num2items(num, len(features))
		for _, item := range items {
			words = append(words, features[item])
		}
		overrideCandidates = append(overrideCandidates, filepath.Join(overridesPath, fmt.Sprintf("%s.yaml", strings.Join(words, "-"))))
	}
	fmt.Fprintln(os.Stderr, "Override candidates:")
	fmt.Fprintln(os.Stderr, strings.Join(overrideCandidates, "\n"))
	overrides := make([]string, 0)
	for _, overrideCandidate := range overrideCandidates {
		if _, err := os.Stat(overrideCandidate); err == nil {
			overrides = append(overrides, fmt.Sprintf("--values %s", overrideCandidate))
		}
	}
	fmt.Println(strings.Join(overrides, " "))
}
