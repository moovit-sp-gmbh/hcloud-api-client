#! /usr/bin/env bash
cat > cfg.yaml << 'EOF'
package: high5Client
output: internal/clients/high5Client/high5Client.gen.go
generate:
  models: true
  client: true
EOF
oapi-codegen -config cfg.yaml <(curl https://eu-dev-1.app.helmut.cloud/api/high5/v1/docs/spec.json | ./patcher) || rm -f cfg.yaml && rm -f cfg.yaml
