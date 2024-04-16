# Install
```
helm plugin install https://github.com/kubepro/helm-plugin-osh.git
```

# Uninstall
```
helm plugin uninstall osh
```

# Usage

## get-values-overrides
```
helm osh get-values-overrides [-d/--download] [-p/--path <lookup_path>] [-u/--url <base_url>] <chart> <feature-1> <feature-2> ...
```

- `-d/--download` download overrides if not found locally in the path
- `-u/--url <base_url>` the base url from where the plugin tries to download overrides. The resulting url is `<base_url>/<chart>/values_overrides/<override_candidate>`
- `-p/--path <lookup_path>` the path where the plugin tries to find override files and puts downloaded override files. The resulting path is `<lookup_path>/<chart>/values_overrides/<override_candidate>`
- `<chart>` the chart name for which the plugin tries to find overrides
- `<feature-X>` features in the beginning of the list are of higher priority

If there are 3 features in the list, then the order of override candidates will be like follows:
```
<lookup_path>/<chart>/values_overrides/<feature-3>.yaml
<lookup_path>/<chart>/values_overrides/<feature-2>.yaml
<lookup_path>/<chart>/values_overrides/<feature-2>-<feature-3>.yaml
<lookup_path>/<chart>/values_overrides/<feature-1>.yaml
<lookup_path>/<chart>/values_overrides/<feature-1>-<feature-3>.yaml
<lookup_path>/<chart>/values_overrides/<feature-1>-<feature-2>.yaml
<lookup_path>/<chart>/values_overrides/<feature-1>-<feature-2>-<feature-3>.yaml
```

## wait-for-pods
KUBECONFIG must be set properly.
```
helm osh wait-for-pods <namespace> [timeout]
```