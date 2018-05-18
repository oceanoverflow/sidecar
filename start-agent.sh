#!/bin/bash

ETCD_HOST=etcd
ETCD_PORT=2379
ETCD_URL=http://$ETCD_HOST:$ETCD_PORT

echo ETCD_URL = $ETCD_URL

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  sidecar consumer --etcd=$ETCD_URL 
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  sidecar provider small --etcd=$ETCD_URL 
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  sidecar provider medium --etcd=$ETCD_URL  
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  sidecar provider large --etcd=$ETCD_URL  
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
