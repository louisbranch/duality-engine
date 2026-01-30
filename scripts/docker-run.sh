#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)
repo_root=$(cd "${script_dir}/.." && pwd)
image_tag=${1:-duality-engine:dev}
data_dir="${repo_root}/data"
user_id=$(id -u)
group_id=$(id -g)

mkdir -p "${data_dir}"

docker build -t "${image_tag}" "${repo_root}"

exec docker run \
	-p 127.0.0.1:8081:8081 \
	-v "${data_dir}:/data" \
	--user "${user_id}:${group_id}" \
	-e DUALITY_DB_PATH=/data/duality.db \
	-e DUALITY_GRPC_ADDR=127.0.0.1:8080 \
	-e DUALITY_MCP_ALLOWED_HOSTS=localhost \
	"${image_tag}"
