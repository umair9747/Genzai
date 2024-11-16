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
    go build
    ./genzai -api &
    GO_PID=$!
}

# Function to start the Streamlit UI
start_streamlit_ui() {
    echo "Starting Streamlit UI..."
    python3 -m streamlit run ./Genzai-UI/ui-main.py &
    STREAMLIT_PID=$!
}

start_go_api
sleep 2
start_streamlit_ui
wait $GO_PID $STREAMLIT_PID

# Cleanup in case processes exit naturally
cleanup
