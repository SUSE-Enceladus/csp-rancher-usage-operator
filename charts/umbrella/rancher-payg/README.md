# How to use this chart
1. Run `helm dependency update .` in this chart to download/update the dependent charts.

2. Identify the appropriate subchart values settings and create an appropriate override values YAML file.

3. Install the chart using a command like the following. Keep the release name as rancher and namespace as cattle-system always.

```console
$ helm upgrade rancher . --namespace cattle-system --create-namespace --install --values ~/overrides.yaml
```
