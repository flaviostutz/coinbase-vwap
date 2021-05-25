#!/bin/bash

echo "Starting coinbase-vwap..."

coinbase-vwap \
  --loglevel="$LOG_LEVEL" \
  --coinbase-ws-url="$COINBASE_WS_URL" \
  --kafka-brokers="$KAFKA_BROKERS" \
  --product-ids="$PRODUCT_IDS"
  

