#!/bin/bash
# server/tests/test_config.sh

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
YELLOW='\033[1;33m'

# Test ports
TEST_PORT_1=4001
TEST_PORT_2=4002
TEST_PORT_3=4003

verify_files() {
    CONFIG_PATH="/home/shamz/Dev/projects/swing-society-website/config.json"
    ENV_PATH="/home/shamz/Dev/projects/swing-society-website/.env"

    echo -e "${YELLOW}DEBUG: Checking if config.json exists at $CONFIG_PATH${NC}"

    if [ ! -f "$CONFIG_PATH" ]; then
        echo -e "${RED}Error: config.json NOT FOUND at $CONFIG_PATH${NC}"
        ls -lah "$(dirname "$CONFIG_PATH")"
        exit 1
    fi

    echo -e "${GREEN}✓ config.json file found at $CONFIG_PATH${NC}"

    echo -e "${YELLOW}DEBUG: Checking if .env exists at $ENV_PATH${NC}"

    if [ ! -f "$ENV_PATH" ]; then
        echo -e "${RED}Error: .env file NOT FOUND at $ENV_PATH${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ .env file found at $ENV_PATH${NC}"
}


# Function to wait for server
wait_for_server() {
    local port=$1
    local retries=10
    while [ $retries -gt 0 ]; do
        if curl -s "http://localhost:$port/health" > /dev/null; then
            return 0
        fi
        sleep 1
        retries=$((retries-1))
    done
    return 1
}

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up processes...${NC}"
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
    fi
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

echo -e "${GREEN}Starting Configuration Tests${NC}\n"

# Verify required files exist
verify_files

# 1. Run Go package tests
echo -e "${GREEN}Running Go package tests for config...${NC}"
go test ./server/internal/config -v -cover

# 2. Test default configuration
echo -e "DEBUG: Current working directory is $(pwd)"

echo -e "\n${GREEN}Testing default configuration...${NC}"
PORT=$TEST_PORT_1 go run server/main.go &
SERVER_PID=$!

if wait_for_server $TEST_PORT_1; then
    echo -e "${GREEN}✓ Server started with default configuration${NC}"
    response=$(curl -s -i "http://localhost:$TEST_PORT_1/health")
    echo "Health check response:"
    echo "$response"
else
    echo -e "${RED}✗ Server failed to start with default configuration${NC}"
fi

kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

# 3. Test custom port configuration
echo -e "\n${GREEN}Testing custom port configuration...${NC}"
PORT=$TEST_PORT_2 go run server/main.go &
SERVER_PID=$!

if wait_for_server $TEST_PORT_2; then
    echo -e "${GREEN}✓ Server started on custom port${NC}"
    response=$(curl -s -i "http://localhost:$TEST_PORT_2/health")
    echo "Health check response:"
    echo "$response"
else
    echo -e "${RED}✗ Server failed to start on custom port${NC}"
fi

kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

# 4. Test production environment
echo -e "\n${GREEN}Testing production environment...${NC}"
ENVIRONMENT=production PORT=$TEST_PORT_3 go run server/main.go &
SERVER_PID=$!

if wait_for_server $TEST_PORT_3; then
    echo -e "${GREEN}✓ Testing cache headers${NC}"
    cache_headers=$(curl -s -I "http://localhost:$TEST_PORT_3/static/css/main.css" | grep -i "cache-control")
    echo "Cache-Control headers:"
    echo "$cache_headers"

    security_headers=$(curl -s -I "http://localhost:$TEST_PORT_3/health" | grep -i "content-security-policy")
    echo "Security headers:"
    echo "$security_headers"
else
    echo -e "${RED}✗ Server failed to start in production mode${NC}"
fi

echo -e "\n${GREEN}Configuration tests completed!${NC}"