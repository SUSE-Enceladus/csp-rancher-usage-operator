{{- define "csp-usage-operator.labels" -}}
app: rancher-csp-usage-operator
{{- end }}

# Namespace where the CSP Billing Adapter is deployed
{{- define "csp-usage-operator.cspBillingAdapterNamespace" -}}
csp-billing-adapter-system
{{- end }}

# Name of the configmap that has the CSP configuration information
{{- define "csp-usage-operator.cspBillingAdapterConfigMap" -}}
csp-config
{{- end }}

# Number of days since last billed before notifying user that we are unable
# to bill them and that they must fix any ongoing issues related to billing
# in order to maintain supportability.
{{- define "csp-usage-operator.cspBillingNoBillThreshold" -}}
30
{{- end }}

{{- define "csp-usage-operator.outputNotification" -}}
csp-compliance
{{- end }}

{{- define "csp-usage-operator.hostnameSetting" -}}
server-url
{{- end }}

{{- define "csp-usage-operator.versionSetting" -}}
server-version
{{- end }}
