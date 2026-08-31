#! /usr/bin/env bash
#
# Regenerate one OpenAPI client.
#
# Usage: scripts/generate.sh <service> <package> <api-path>
#
#   service   name for the vendored spec, and the shorthand codegen knows
#   package   Go package to generate, e.g. cosmoClient
#   api-path  path segment under /api on the dev cluster
#
# The patched spec is vendored under specs/ as well as piped into
# oapi-codegen: the specs declare no operationIds and oapi-codegen drops their
# summaries, so specs/ is the only place the human-readable text and the tags
# survive. cmd/codegen reads them when it names a generated command.
set -euo pipefail

service=$1
package=$2
api_path=$3

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"

if [[ ! -x ./patcher ]]; then
	echo "generate: ./patcher not built — run 'make build' first" >&2
	exit 1
fi

cluster=${HCLOUD_SPEC_CLUSTER:-https://eu-dev-1.app.helmut.cloud}

# A shared cfg.yaml in the repo root would stop the per-service targets running
# in parallel, so keep it in a private temp dir.
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

cat > "$tmp/cfg.yaml" <<CFG
package: $package
output: internal/clients/$package/$package.gen.go
generate:
  models: true
  client: true
CFG

mkdir -p specs
curl --fail --silent --show-error "$cluster/api/$api_path/v1/docs/spec.json" \
	| ./patcher > "specs/$service.json"

oapi-codegen -config "$tmp/cfg.yaml" "specs/$service.json"
