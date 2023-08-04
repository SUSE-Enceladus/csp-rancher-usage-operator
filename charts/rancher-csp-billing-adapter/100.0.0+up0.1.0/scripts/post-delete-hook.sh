#!/bin/bash

set -e

rancher_namespace="${RANCHER_NAMESPACE}"

if [[ -z ${rancher_namespace} ]]; then
  echo "No rancher namespace is provided."
  exit 1
fi

failed=()

echo "Uninstalling Rancher resources in the following namespace: ${rancher_namespace}"

if [[ ! $(helm uninstall rancher -n "${rancher_namespace}") ]]; then
  failed=("${failed[@]}" "rancher")
fi

if [[ ${#failed[@]} -ne 0 ]]; then
  echo "Failed to uninstall the following apps:" "${failed[@]}"
else
  echo "Rancher uninstalled successfully"
fi
