#!/bin/sh

set -eux
cd $(dirname $0)

TEMPDIR=$(mktemp -d)
trap "rm -rf $TEMPDIR" EXIT

GO_PACKAGE_NAME=$(awk '/^module/{print $2}' go.mod)/azure-sdk-for-go
API_SPEC_PATH=$TEMPDIR/specification/consumption/resource-manager

git clone --depth=1 https://github.com/Azure/azure-rest-api-specs.git $TEMPDIR

cat <<'EOF' >>$API_SPEC_PATH/readme.go.md

``` yaml $(tag) == 'package-2019-10' & $(go)
output-folder: $(go-sdk-folder)/services/$(namespace)/mgmt/2019-10-01/$(namespace)
package-name: $(go-package-name)/services/$(namespace)/mgmt/2019-10-01/$(namespace)
```
EOF

rm -rf azure-sdk-for-go

autorest \
    --go \
    --tag=package-2019-10 \
    --go-sdk-folder=./azure-sdk-for-go \
    --go-package-name=$GO_PACKAGE_NAME \
    $API_SPEC_PATH

gofmt -w azure-sdk-for-go
