#!/bin/bash

set -e

# Set console colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}       Project EduNexus - One-Click Start       ${NC}"
echo -e "${BLUE}================================================${NC}"

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}Error: 'go' command not found. Please install Go environment first.${NC}"
    exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
    echo -e "${RED}Error: 'npm' command not found. Please install Node.js and npm first.${NC}"
    exit 1
fi

# Check project directory
if [ ! -d "edunexus" ]; then
    if [ ! -f "wails.json" ]; then
        echo -e "${RED}Error: Please run this script in the project root or edunexus directory.${NC}"
        exit 1
    fi
else
    cd edunexus
fi

echo -e "\n${GREEN}[1/2] Building Frontend...${NC}"
cd frontend
npm install
npm run build
cd ..

echo -e "\n${GREEN}[2/2] Starting EduNexus Client...${NC}"
echo -e "${YELLOW}Starting with go run -tags dev . (Press Ctrl+C to stop)${NC}"

# Run wails project
go run -tags dev .
