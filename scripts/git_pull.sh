#!/bin/bash

# Configuration variables
REMOTE_USER="root"
REMOTE_DIR="/root/mp2-distributed-systems"

# List of remote hosts
REMOTE_HOSTS=(
    "fa24-cs425-9201.cs.illinois.edu"
    "fa24-cs425-9202.cs.illinois.edu"
    "fa24-cs425-9203.cs.illinois.edu"
    "fa24-cs425-9204.cs.illinois.edu"
    "fa24-cs425-9205.cs.illinois.edu"
    "fa24-cs425-9206.cs.illinois.edu"
    "fa24-cs425-9207.cs.illinois.edu"
    "fa24-cs425-9208.cs.illinois.edu"
    "fa24-cs425-9209.cs.illinois.edu"
    "fa24-cs425-9210.cs.illinois.edu"
)

# Function to perform git pull on a remote host
perform_git_pull() {
    local host=$1
    echo "Connecting to $host..."
    ssh -o StrictHostKeyChecking=no $REMOTE_USER@$host << EOF
        echo "Connected to $host"
        cd $REMOTE_DIR
        echo "Changed directory to $REMOTE_DIR"
        echo "Performing git pull..."
        git fetch
        git checkout main
        git pull
        echo "Git pull completed on $host"
        exit
EOF
    echo "Disconnected from $host"
    echo "------------------------"
}

# Main execution
for host in "${REMOTE_HOSTS[@]}"
do
    perform_git_pull $host
done

echo "Script execution completed for all hosts"
