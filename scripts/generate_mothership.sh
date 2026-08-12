#! /usr/bin/env bash
cat > cfg.yaml << 'EOF'
package: mothershipClient
output: internal/clients/mothershipClient/mothershipClient.gen.go
generate:
  models: true
  client: true
EOF
oapi-codegen -config cfg.yaml <(curl https://eu-dev-1.app.helmut.cloud/api/mothership/v1/docs/spec.json | ./patcher) || rm -f cfg.yaml && rm -f cfg.yaml
