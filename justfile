# Default just command
default:
    @just --list

export:
    #!/usr/bin/env bash
    set -euxo pipefail
    dagger call build-env --source=. export --path=nobuffer.tgz
    docker import nobuffer.tgz
    rm -f nobuffer.tgz

# Run tests using Dagger
test-default:
    dagger call test --source=.

# Run tests using Dagger
test:
    dagger call test --source=. --image-name=alpine

help:
    dagger call test --source=. --help

# Build and run the exported image
run IMAGE_NAME="nobuffer" CONTAINER_NAME="nobuffer-instance":
    docker run --rm --tty --interactive --name {{CONTAINER_NAME}} {{IMAGE_NAME}} sleep infinity

# Execute a command in the running container
exec CONTAINER_NAME="nobuffer-instance" COMMAND="/bin/sh":
    docker exec --tty --interactive {{CONTAINER_NAME}} {{COMMAND}}

# Build and publish image with OCI labels
publish REGISTRY_URL="ttl.sh/my-nobuffer:24h":
    dagger call build-and-publish --source=. --image-name=pandoc/core --registry-url={{REGISTRY_URL}}

# Fetch the published image and run a container from it
fetch-and-run IMAGE_URL="ttl.sh/my-nobuffer:24h" CONTAINER_NAME="nobuffer-instance":
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Pull the latest image
    docker pull {{IMAGE_URL}}
    
    # Check if the container already exists
    if docker ps -a --format '{{ '{{' }}.Names{{ '}}' }}' | grep -q '^{{CONTAINER_NAME}}$'; then
        echo "Container {{CONTAINER_NAME}} already exists. Removing it..."
        docker rm --force {{CONTAINER_NAME}}
    fi
    
    # Run the new container
    docker run --name {{CONTAINER_NAME}} --detach {{IMAGE_URL}} --entrypoint=bash -c "sleep infinity"
    
    echo "Container {{CONTAINER_NAME}} is now running. Use 'just exec {{CONTAINER_NAME}}' to interact with it."
    echo "To stop and remove the container, use 'docker stop {{CONTAINER_NAME}} && docker rm {{CONTAINER_NAME}}'"
