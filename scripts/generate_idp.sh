#! /usr/bin/env bash
cat > cfg.yaml << 'EOF'
package: idpClient
output: internal/clients/idpClient/idpClient.gen.go
generate:
  models: true
  client: true
EOF
oapi-codegen -config cfg.yaml <(curl https://eu-dev-1.app.helmut.cloud/api/account/v1/docs/spec.yaml) || rm -f cfg.yaml && rm -f cfg.yaml
