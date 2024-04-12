# Install
```
helm plugin install https://github.com/kubepro/helm-get-values-overrides.git
```

# Uninstall
```
helm plugin uninstall get-values-overrides
```

# Usage
```
helm get-values-overrides <lookup_path> <feature-1> <feature-2> ...
```

- `lookup_path` the path where the plugin tries to find values override files.
- `feature-X` features in the beginning of the list are of higher priority.

If there are 3 features in the list, then the order of override candidates will be like follows:
```
<lookup_path>/<feature-3>.yaml
<lookup_path>/<feature-2>.yaml
<lookup_path>/<feature-2>-<feature-3>.yaml
<lookup_path>/<feature-1>.yaml
<lookup_path>/<feature-1>-<feature-3>.yaml
<lookup_path>/<feature-1>-<feature-2>.yaml
<lookup_path>/<feature-1>-<feature-2>-<feature-3>.yaml
```
