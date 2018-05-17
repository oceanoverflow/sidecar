#!/bin/bash

ETCD_HOST=$(ip addr show docker0 | grep 'inet\b' | awk '{print $2}' | cut -d '/' -f 1)
ETCD_PORT=2379
ETCD_URL=http://$ETCD_HOST:$ETCD_PORT

echo ETCD_URL = $ETCD_URL

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  /home/sidecar consumer --etcd=$ETCD_URL 
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  /home/sidecar provider small --etcd=$ETCD_URL 
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  /home/sidecar provider medium --etcd=$ETCD_URL  
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  /home/sidecar provider large --etcd=$ETCD_URL  
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
