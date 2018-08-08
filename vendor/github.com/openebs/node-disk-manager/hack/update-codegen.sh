#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CODEGEN_PKG=${GOPATH}/src/github.com/weaveworks/scope

${CODEGEN_PKG}/vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/weaveworks/scope/vendor/github.com/openebs/node-disk-manager/pkg/client github.com/weaveworks/scope/vendor/github.com/openebs/node-disk-manager/pkg/apis \
  openebs.io:v1alpha1
