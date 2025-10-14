#!/bin/bash

# Get yarn on the path
export NVM_DIR="$HOME/.nvm" && \
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && \

# install dependencies and generate prisma client
go mod tidy
yarn install
nx generate-client db

# start supervisord
mkdir -p /var/log/supervisor
supervisord -c ./deploy/engine/local/config/supervisor/supervisord.conf

# tail the supervisor logs
tail -f /var/log/supervisor/*.log