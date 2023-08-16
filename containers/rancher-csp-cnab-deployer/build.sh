#!/bin/sh

PROJECT_DIR=$(dirname -- $(readlink -e -- ${BASH_SOURCE[0]}))
: ${DEPLOYER_IMAGE_REPO:="yeey/rancher-csp-cnab-deployer"}
: ${DEPLOYER_IMAGE_TAG:="latest"}

set -xe

pushd $PROJECT_DIR
rm -rf app/charts/*
for chart in ../../charts/*
do
	if [ -z "${chart##*cnab-deployer*}" ]
	then
		# don't copy the deployer chart
		continue
	fi
	cp -r $chart app/charts
done
docker build -t ${DEPLOYER_IMAGE_REPO}:${DEPLOYER_IMAGE_TAG} .

rm -rf app/charts/*

docker push ${DEPLOYER_IMAGE_REPO}:${DEPLOYER_IMAGE_TAG}
popd
