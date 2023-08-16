# Build Container

Before building the `rancher-csp-cnab-deployer` container, make sure you have
an account with docker hub and create a repo called
`rancher-csp-cnab-deployer`.

To build the `rancher-csp-cnab-deployer` container:

1. Set the environment variable `DEPLOYER_IMAGE_REPO` environment variable to
   the `rancher-csp-cnab-deployer` repo in your docker hub.
2. Optionally set the environment variable `DEPLOYER_IMAGE_TAG`. If unset,
   `latest` will be used by default.
3. Run `build.sh` to build the container and push it to your docker hub repo.

For example:

```console
export DEPLOYER_IMAGE_REPO="yeey/rancher-csp-cnab-deployer"
./build.sh
```

# Deploy The Container

After building it, you can use the `rancher-csp-cnab-deployer` chart to deploy
Rancher, usage operator CRD, usage operator, and billing adapter.

For example,

```console
cd ../../charts
helm install rancher-csp-cnab-deployer rancher-csp-cnab-deployer -n cattle-rancher-csp-deployer-system --create-namespace --atomic
```

# Unstall

To unstall the the deployer:

```console
helm uninstall rancher-csp-cnab-deployer -n cattle-rancher-csp-deployer-system --debug
```
