{{- define "csp-usage-operator.labels" -}}
app: rancher-csp-usage-operator
{{- end }}

{{- define "csp-usage-operator.outputConfigMap" -}}
csp-config
{{- end }}

{{- define "csp-usage-operator.outputNotification" -}}
csp-compliance
{{- end }}

{{- define "csp-usage-operator.cacheSecret" -}}
csp-usage-operator-cache
{{- end }}

{{- define "csp-usage-operator.hostnameSetting" -}}
server-url
{{- end }}

{{- define "csp-usage-operator.versionSetting" -}}
server-version
{{- end }}
