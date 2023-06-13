{{/*
Expand the name of the chart.
*/}}
{{- define "csp-rancher-usage-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "csp-rancher-usage-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "csp-rancher-usage-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "csp-rancher-usage-operator.labels" -}}
helm.sh/chart: {{ include "csp-rancher-usage-operator.chart" . }}
{{ include "csp-rancher-usage-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "csp-rancher-usage-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "csp-rancher-usage-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "csp-rancher-usage-operator.serviceAccountName" -}}
{{- default (include "csp-rancher-usage-operator.fullname" .) .Values.serviceAccount.name }}
{{- end }}

# Namespace where the CSP Billing Adapter is deployed
{{- define "csp-usage-operator.cspBillingAdapterNamespace" -}}
cattle-csp-billing-adapter-system
{{- end }}

# Name of the configmap that has the CSP configuration information
{{- define "csp-usage-operator.cspBillingAdapterConfigMap" -}}
csp-config
{{- end }}

# Number of days since last billed before notifying user that we are unable
# to bill them and that they must fix any ongoing issues related to billing
# in order to maintain supportability.
{{- define "csp-usage-operator.cspBillingNoBillThreshold" -}}
45
{{- end }}

# Interval in whihc Usage Operator runs.
# Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
{{- define "csp-usage-operator.cspBillingManagerInterval" -}}
30s
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
