#!/bin/bash

TIMEOUT=30

echo "Waiting for RabbitMQ to be up"
i=0
until curl -s "${RABBITMQ_ENDPOINT}/api" > /dev/null; do
    i=$((i + 1))
    if [ $i -eq $TIMEOUT ]; then
        echo
        echo "Timeout while waiting for RabbitMQ to be up"
        exit 1
    fi
    printf "."
    sleep 2
done
