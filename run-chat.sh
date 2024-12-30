#!/bin/bash

# Default values
NATS_URL="nats://localhost:4222"
NAME="anonymous"
CHANNEL="general"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -nats)
      NATS_URL="$2"
      shift 2
      ;;
    -name)
      NAME="$2"
      shift 2
      ;;
    -channel)
      CHANNEL="$2"
      shift 2
      ;;
    *)
      echo "Unknown parameter: $1"
      echo "Usage: $0 [-nats <NATS URL>] [-name <username>] [-channel <channel name>]"
      exit 1
      ;;
  esac
done

# Check if image exists
if [[ "$(docker images -q chat-cli 2> /dev/null)" == "" ]]; then
  echo "Building chat-cli image..."
  docker build -t chat-cli .
fi

# Run the container
echo "Connecting to $NATS_URL as $NAME in channel $CHANNEL"
docker run -it --network host chat-cli \
  -nats "$NATS_URL" \
  -name "$NAME" \
  -channel "$CHANNEL" 