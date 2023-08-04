#!/bin/bash

set -e

rancher_namespace="${RANCHER_NAMESPACE}"
rancher_hostname=${RANCHER_HOSTNAME}

if [[ -z ${rancher_namespace} ]]; then
  echo "No rancher namespace is provided."
  exit 1
fi

if [[ -z ${rancher_hostname} ]]; then
  echo "No rancher hostname is provided."
  exit 1
fi

failed=()

echo "installing Rancher resources in the following namespace: ${rancher_namespace} ${rancher_hostname}"

helm repo add rancher-latest https://releases.rancher.com/server-charts/latest

if [[ ! $(helm upgrade --install rancher rancher-latest/rancher --namespace ${rancher_namespace} --create-namespace --set installCRDs=true --set hostname=${rancher_hostname} --set replicas=1 --set extraEnv[0].name=CATTLE_PROMETHEUS_METRICS --set-string extraEnv[0].value="true" --set extraEnv[1].name=CATTLE_SERVER_URL --set extraEnv[1].value=https://${rancher_hostname}) ]]; then
  failed=("${failed[@]}" "rancher")
fi

kubectl -n cattle-system rollout status deploy/rancher

echo "------ Summary ------"
if [[ ${#failed[@]} -ne 0 ]]; then
  echo "Failed to install the following apps:" "${failed[@]}"
else
  echo "Rancher installed successfully."
fi
