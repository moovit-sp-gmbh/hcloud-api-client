#! /usr/bin/env bash
cat > cfg.yaml << 'EOF'
package: cosmoClient
output: internal/clients/cosmoClient/cosmoClient.gen.go
generate:
  models: true
  client: true
EOF
oapi-codegen -config cfg.yaml <(curl https://eu-dev-1.app.helmut.cloud/api/cosmo/v1/docs/spec.json | ./patcher) || rm -f cfg.yaml && rm -f cfg.yaml
