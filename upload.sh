#!/usr/bin/env bash

# Check if help is requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 <command> [args...]"
    echo "Commands:"
    echo "  upload <file1> [file2...]  Upload files"
    echo "  list                         List files"
    echo "  download <URL>              Download file from URL"
    echo "  delete <URL or file path>   Delete file"
    exit 0
fi

# Check if there are enough arguments
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <command> [args...]"
    echo "Use -h or --help for more information"
    exit 1
fi

# First argument is the command
command=$1
shift

# Base URL for operations
base_url="http://localhost:8080"

case $command in
    upload)
        if [ "$#" -lt 1 ]; then
            echo "Usage: $0 upload <file1> [file2...]"
            exit 1
        fi
        for file; do
            if [ ! -f "$file" ]; then
                echo "Error: File '$file' does not exist."
                continue
            fi
            echo "Uploading $file..."
            response=$(curl -s -F "file=@$file" "$base_url/upload")
            if jq -e '.url' >/dev/null 2>&1 <<< "$response"; then
                url=$(jq -r '.url' <<< "$response")
                echo "File uploaded: $url"
                wl-copy <<< "$url"
            else
                echo "Error uploading file: $response"
            fi
        done
        ;;
    download)
        if [ "$#" -ne 1 ]; then
            echo "Usage: $0 download <URL>"
            exit 1
        fi
        curl --fail -O "$1"
        ;;
    delete)
        if [ "$#" -ne 1 ]; then
            echo "Usage: $0 delete <URL or file path>"
            exit 1
        fi
        curl -X DELETE "${1/#\//${base_url}/}"
        ;;
    list)
        curl --silent "$base_url/files" | jq .
        ;;
    *)
        echo "Invalid command: $command"
        echo "Usage: $0 <command> [args...]"
        echo "Commands:"
        echo "  upload <file1> [file2...]  Upload files"
        echo "  list                         List files"
        echo "  download <URL>              Download file from URL"
        echo "  delete <URL or file path>   Delete file"
        exit 1
        ;;
esac