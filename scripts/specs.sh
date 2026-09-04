#! /usr/bin/env bash
#
# Vendor the patched OpenAPI specs under specs/ without regenerating any client.
#
# The specs are the only source of the operation summaries, descriptions,
# parameter descriptions, tags and security requirements: they declare no
# operationIds and oapi-codegen drops all of that, so the generated clients
# carry none of it. cmd/codegen reads specs/ to prefill what it asks about.
#
# Regenerating a client rewrites tens of thousands of lines, so this is split
# out from `make generate` to make refreshing the metadata cheap.
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"

if [[ ! -x ./patcher ]]; then
	echo "specs: ./patcher not built — run 'make build' first" >&2
	exit 1
fi

cluster=${HCLOUD_SPEC_CLUSTER:-https://eu-dev-1.app.helmut.cloud}

# service:api-path — the api path differs from the service name for idp.
services=(
	"idp:account"
	"high5:high5"
	"mothership:mothership"
	"cosmo:cosmo"
	"dali:dali"
)

mkdir -p specs
for entry in "${services[@]}"; do
	service=${entry%%:*}
	api_path=${entry##*:}

	curl --fail --silent --show-error "$cluster/api/$api_path/v1/docs/spec.json" \
		| ./patcher > "specs/$service.json"
	printf 'specs/%s.json\n' "$service"
done
