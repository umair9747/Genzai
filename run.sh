#!/bin/bash

# Function to cleanup and exit
cleanup() {
    echo "Stopping all services..."
    [[ -n $GO_PID ]] && kill $GO_PID
    [[ -n $STREAMLIT_PID ]] && kill $STREAMLIT_PID
    exit 0
}

# Trap Ctrl+C (SIGINT) and call cleanup
trap cleanup SIGINT

# Function to start the Go API
start_go_api() {
    echo "Building and starting Go API..."
    cd Genzai-Tool || exit
    go build
    ./genzai -api &
    GO_PID=$!
    cd ..
}

# Function to start the Streamlit UI
start_streamlit_ui() {
    echo "Starting Streamlit UI..."
    python3 -m streamlit run Genzai-UI/ui-main.py &
    STREAMLIT_PID=$!
}

# Parse command-line arguments
if [[ "$1" == "-api" && "$2" == "-ui" ]]; then
    start_go_api
    sleep 2
    start_streamlit_ui
    wait $GO_PID $STREAMLIT_PID
elif [[ "$1" =~ ^http:// ]]; then
    echo "Running ./genzai on target $1..."
    cd Genzai-Tool || exit
    ./genzai "$1"
else
    echo "Usage:"
    echo "./run.sh -api -ui         # Start API and UI"
    echo "./run.sh http://ip:port   # Run on target without UI"
    exit 1
fi

# Cleanup in case processes exit naturally
cleanup