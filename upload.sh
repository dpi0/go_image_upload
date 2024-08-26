#!/usr/bin/env bash

# Check if there are enough arguments for the commands that require them
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <file1> [file2...]"
    echo "Or for other commands: $0 <command> [args...]"
    echo "Commands: list, download, delete"
    exit 1
fi

# First argument can be a command or a file
first_arg=$1

# Check if the first argument is a command
if [[ "$first_arg" == "download" || "$first_arg" == "delete" || "$first_arg" == "list" ]]; then
    command=$first_arg
    shift
else
    command="upload"
fi

# Base URL for operations
base_url="http://localhost:8081"

case $command in
    upload)
        # Shift to remove the "upload" keyword itself from the list
        shift
        for file in "$@"; do
            echo "Uploading $file..."
            response=$(curl -s -F "file=@$file" "$base_url/upload")
            url=$(echo $response | jq -r '.url')
            echo "File uploaded. $url"
            echo $url | wl-copy
        done
        ;;
    download)
        if [ "$#" -ne 1 ]; then
            echo "Usage: $0 download <URL>"
            exit 1
        fi
        url=$1
        echo "Downloading from $url..."
        curl -O "$url"
        ;;
    delete)
        if [ "$#" -ne 1 ]; then
            echo "Usage: $0 delete <URL or file path>"
            exit 1
        fi
        url=$1
        if [[ $url != http* ]]; then
            url="$base_url$url"
        fi
        echo "Deleting $url..."
        curl -X DELETE "$url"
        ;;
    list)
        echo "Listing files from $base_url/files..."
        curl --silent "$base_url/files" | jq .
        ;;
    *)
        echo "Invalid command: $command"
        echo "Usage: $0 <file1> [file2...]"
        echo "Or for other commands: $0 <command> [args...]"
        echo "Commands: list, download, delete"
        exit 1
        ;;
esac
