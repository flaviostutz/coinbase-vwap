#!/bin/bash

echo "Starting coinbase-vwap..."

coinbase-vwap \
  --loglevel="$LOG_LEVEL" \
  --coinbase-ws-url="$COINBASE_WS_URL"

