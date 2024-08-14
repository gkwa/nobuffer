# Default just command
default:
    @just --list

# Export the image from Dagger and import it into Docker
export-import IMAGE_NAME="nobuffer":
    #!/usr/bin/env bash
    set -euxo pipefail
    dagger export -o image.tar
    docker import image.tar {{IMAGE_NAME}}
    rm image.tar

# Run tests using Dagger
test-default:
    dagger call test --source=.

# Run tests using Dagger
test:
    dagger call test \
        --source=. \
        --lua-version=5.4 \
        --image-name=alpine \
        --image-version=3.18

help:
    dagger call test --source=. --help

# Build, export, import, and test in one command
all: export-import test

# Clean up build artifacts
clean:
    rm -f nobuffer image.tar

# Build and run the exported image
run IMAGE_NAME="nobuffer" CONTAINER_NAME="nobuffer-instance":
    docker run --rm --tty --interactive --name {{CONTAINER_NAME}} {{IMAGE_NAME}} sleep infinity

# Execute a command in the running container
exec CONTAINER_NAME="nobuffer-instance" COMMAND="/bin/sh":
    docker exec --tty --interactive {{CONTAINER_NAME}} {{COMMAND}}

# Build and publish image with OCI labels
publish REGISTRY_URL="ttl.sh/my-nobuffer:24h":
    dagger call build-and-publish --source=. --lua-version=5.4 --image-name=alpine --image-version=3.18 --registry-url={{REGISTRY_URL}}

# Fetch the published image and run a container from it
fetch-and-run IMAGE_URL="ttl.sh/my-nobuffer:latest" CONTAINER_NAME="nobuffer-instance":
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
    docker run --name {{CONTAINER_NAME}} --detach {{IMAGE_URL}} sleep infinity
    
    echo "Container {{CONTAINER_NAME}} is now running. Use 'just exec {{CONTAINER_NAME}}' to interact with it."
    echo "To stop and remove the container, use 'docker stop {{CONTAINER_NAME}} && docker rm {{CONTAINER_NAME}}'"
