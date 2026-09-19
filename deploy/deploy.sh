#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
set -a; source .env; set +a

echo "[$(date -Iseconds)] deploy $IMAGE:$IMAGE_TAG"

docker compose pull app

mkdir -p /opt/backup
docker compose exec -T postgres pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" \
  | gzip > "/opt/backup/${STACK_NAME}-$(date +%F-%H%M).sql.gz"
echo "[$(date -Iseconds)] backup selesai"

docker compose run --rm migrate
echo "[$(date -Iseconds)] migration selesai"

docker compose up -d app

for i in $(seq 1 30); do
  status=$(docker inspect -f '{{.State.Health.Status}}' "${STACK_NAME}-app" 2>/dev/null || echo none)
  if [ "$status" = "healthy" ]; then
    echo "[$(date -Iseconds)] deploy sukses"
    docker image prune -f
    exit 0
  fi
  sleep 2
done

echo "[$(date -Iseconds)] GAGAL: container tidak healthy dalam 60 detik" >&2
docker compose logs --tail=50 app >&2
exit 1