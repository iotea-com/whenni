#!/bin/bash

# This script's context is the deploy/services/local folder

# Exit immediately if a command exits with a non-zero status
set -e

# Path to the Docker Compose file
DOCKER_COMPOSE_FILE="docker-compose.yaml"

# Profile to run, default to 'all' if not provided
PROFILE=${2:-all}

# Profiles for docker compose
COMPOSE_PROFILES=""
valid_profiles=("source" "processing" "action" "conditional" "custom" "all")

# Check if the profile is valid
if [[ " ${valid_profiles[@]} " =~ " ${PROFILE} " ]]; then
  if [ "$PROFILE" = "all" ]; then
    COMPOSE_PROFILES="--profile source --profile processing --profile action --profile conditional --profile custom"
  else
    COMPOSE_PROFILES="--profile $PROFILE"
  fi
else
  echo "Usage: $0 {source|processing|action|conditional|custom|all}"
  exit 1
fi

function is_service_running() {
  local service=$1
  local state=$(docker compose -f "$DOCKER_COMPOSE_FILE" ps -q $service | xargs docker inspect -f '{{.State.Running}}')
  echo "$state"
}

# Check if the Docker Compose file exists
function check_docker_file {
  if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
    echo "Docker Compose file not found: $DOCKER_COMPOSE_FILE"
    exit 1
  fi
}

# Function to start Docker Compose services
function start {
  # Temporarily list services to check if any are available for the selected profiles
  local services=$(docker compose -f "$DOCKER_COMPOSE_FILE" $COMPOSE_PROFILES config --services)

  if [ -z "$services" ]; then
    echo "Warning: No services are configured for the selected profiles: $PROFILE. This may be intentional for future extensions."
    return
  fi

  # Start the services with the selected profiles
  docker compose -f "$DOCKER_COMPOSE_FILE" $COMPOSE_PROFILES up -d

  # Flag to keep track of whether all services are up
  local all_services_up=false

  # Maximum number of retries
  max_retries=30
  retry_count=0
  sleep_interval=1 # 1 second

  # Loop until all services are running or max retries reached
  while [ "$all_services_up" = false ] && [ $retry_count -lt $max_retries ]; do
    # Wait some time initially 
    sleep 1

    all_services_up=true

    for service in $services; do
      local state=$(is_service_running $service)

      if [ "$state" != "true" ]; then
        echo "Service $service is not running. Retry $retry_count/$max_retries."
        all_services_up=false
        break
      fi
    done

    # If all services are up, break the loop
    if [ "$all_services_up" = true ]; then
      echo "All services are running."
      break
    fi

    # Increment retry count and wait before next check
    ((retry_count++))
    sleep $sleep_interval
  done

  # Check final status
  if [ "$all_services_up" = false ]; then
    echo "Some services are still not running after $max_retries retries."
    exit 1
  else
    echo "All services are successfully running."
  fi
}

# Main function to orchestrate the script
function main {
  check_docker_file

  # Sanitizing input
  if [ "$1" != "up" ] && [ "$1" != "stop" ] && [ "$1" != "down" ]; then
      echo "Error: Argument must be 'up', 'stop' or 'down'."
      exit 1
  fi

  # Run the script
  case $1 in
      up)
          echo "Starting services..."
          start
          ;;
      stop)
          echo "Stopping services..."
          docker compose -f $DOCKER_COMPOSE_FILE $COMPOSE_PROFILES stop
          ;;
      down)
          echo "Stopping..."
          docker compose -f $DOCKER_COMPOSE_FILE $COMPOSE_PROFILES down
          ;;
  esac
}

main $1