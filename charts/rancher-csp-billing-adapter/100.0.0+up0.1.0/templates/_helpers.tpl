{{/*
Expand the name of the chart.
*/}}
{{- define "rancher-csp-billing-adapter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rancher-csp-billing-adapter.fullname" -}}
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
{{- define "rancher-csp-billing-adapter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "rancher-csp-billing-adapter.labels" -}}
helm.sh/chart: {{ include "rancher-csp-billing-adapter.chart" . }}
{{ include "rancher-csp-billing-adapter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "rancher-csp-billing-adapter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "rancher-csp-billing-adapter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "rancher-csp-billing-adapter.serviceAccountName" -}}
{{- default (include "rancher-csp-billing-adapter.fullname" .) .Values.serviceAccount.name }}
{{- end }}

# Name of the configmap that has the CSP configuration information
{{- define "csp-billing-adapter.cspBillingAdapterConfigMap" -}}
csp-config
{{- end }}

{{- define "csp-billing-adapter.usageCRDPlural" -}}
cspadapterusagerecords
{{- end }}

{{- define "csp-billing-adapter.usageResource" -}}
rancher-usage-record
{{- end }}

{{- define "csp-billing-adapter.usageAPIGroup" -}}
susecloud.net
{{- end }}

{{- define "csp-billing-adapter.usageAPIVersion" -}}
v1
{{- end }}
