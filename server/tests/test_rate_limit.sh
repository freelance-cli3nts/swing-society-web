#!/bin/bash
for i in {1..10}; do
    echo "Request $i"
    curl -I http://localhost:8080
    sleep 0.1
done
