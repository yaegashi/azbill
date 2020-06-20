#!/bin/sh

set -eux
cd $(dirname $0)

GO_PACKAGE_NAME=$(awk '/^module/{print $2}' go.mod)/azure-sdk-for-go

rm -rf specs azure-sdk-for-go

git clone --depth=1 https://github.com/Azure/azure-rest-api-specs.git specs

cat <<'EOF' >>specs/specification/consumption/resource-manager/readme.go.md

``` yaml $(tag) == 'package-2019-10' & $(go)
output-folder: $(go-sdk-folder)/services/$(namespace)/mgmt/2019-10-01/$(namespace)
package-name: $(go-package-name)/services/$(namespace)/mgmt/2019-10-01/$(namespace)
```
EOF

autorest \
    --go \
    --tag=package-2019-10 \
    --go-sdk-folder=./azure-sdk-for-go \
    --go-package-name=$GO_PACKAGE_NAME \
    specs/specification/consumption/resource-manager

rm -rf specs

gofmt -w azure-sdk-for-go
