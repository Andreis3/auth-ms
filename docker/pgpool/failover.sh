#!/bin/bash
# args: node_id, hostname, port, role
NODE_ID=$1
HOST=$2
PORT=$3
ROLE=$4

echo "$(date) [failover] Node $NODE_ID ($ROLE) caiu; promovendo standby em $HOST:$PORT..."

# roda pg_promote() na réplica que virou primary
PGPASSWORD=admin psql -h "$HOST" -p "$PORT" -U admin -d postgres \
  -c "SELECT pg_promote(wait_seconds => 60, promoting_mode => 'immediate');"
