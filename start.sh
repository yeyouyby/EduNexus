#!/bin/bash

# Do not exit immediately on error to allow for error messages

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
    cd edunexus || {
        echo -e "${RED}Error: Failed to change directory to 'edunexus'.${NC}"
        exit 1
    }
fi

echo -e "\n${GREEN}[1/2] Building Frontend...${NC}"
cd frontend || {
    echo -e "${RED}Error: Failed to change directory to 'frontend'.${NC}"
    exit 1
}

if ! npm install; then
    echo -e "${RED}Error: 'npm install' failed.${NC}"
    exit 1
fi

if ! npm run build; then
    echo -e "${RED}Error: 'npm run build' failed.${NC}"
    exit 1
fi
cd .. || {
    echo -e "${RED}Error: Failed to change directory back from 'frontend'.${NC}"
    exit 1
}

echo -e "\n${GREEN}[2/2] Starting EduNexus Client...${NC}"

# Check if Wails is installed
if command -v wails >/dev/null 2>&1; then
    echo -e "${YELLOW}Wails CLI found. Starting with 'wails dev' (Press Ctrl+C to stop)${NC}"
    wails dev
else
    echo -e "${YELLOW}Wails CLI not found. Starting with 'go run -tags dev .' (Press Ctrl+C to stop)${NC}"
    go run -tags dev .
fi

EXIT_CODE=$?
if [ $EXIT_CODE -ne 0 ] && [ $EXIT_CODE -ne 130 ]; then
    echo -e "\n${RED}Error: The application failed to start or exited with an error (Exit code: $EXIT_CODE).${NC}"
    exit $EXIT_CODE
else
    echo -e "\n${GREEN}Application closed.${NC}"
fi
