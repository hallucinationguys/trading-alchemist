#!/bin/bash

# Build script for Trading Alchemist Docker images

set -e

# Default values
IMAGE_NAME="trading-alchemist"
DOCKERFILE="Dockerfile"
CONTEXT="."
TAG="latest"
REGISTRY=""
PUSH=false
BUILD_ARGS=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -f|--file)
            DOCKERFILE="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        --push)
            PUSH=true
            shift
            ;;
        --build-arg)
            BUILD_ARGS="$BUILD_ARGS --build-arg $2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -n, --name NAME        Image name (default: trading-alchemist)"
            echo "  -t, --tag TAG          Image tag (default: latest)"
            echo "  -f, --file FILE        Dockerfile path (default: Dockerfile)"
            echo "  -r, --registry REG     Registry prefix"
            echo "      --push             Push image after build"
            echo "      --build-arg ARG    Build argument (can be used multiple times)"
            echo "  -h, --help             Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Construct full image name
if [ -n "$REGISTRY" ]; then
    FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME:$TAG"
else
    FULL_IMAGE_NAME="$IMAGE_NAME:$TAG"
fi

echo "Building Docker image..."
echo "Image: $FULL_IMAGE_NAME"
echo "Dockerfile: $DOCKERFILE"
echo "Context: $CONTEXT"

# Build the image
docker build $BUILD_ARGS -f "$DOCKERFILE" -t "$FULL_IMAGE_NAME" "$CONTEXT"

echo "Build completed successfully!"

# Tag with additional tags if this is latest
if [ "$TAG" = "latest" ]; then
    # Try to get version from git tag or use timestamp
    if git describe --tags --exact-match HEAD 2>/dev/null; then
        VERSION=$(git describe --tags --exact-match HEAD)
        VERSION_TAG="${REGISTRY:+$REGISTRY/}$IMAGE_NAME:$VERSION"
        docker tag "$FULL_IMAGE_NAME" "$VERSION_TAG"
        echo "Tagged as: $VERSION_TAG"
    else
        TIMESTAMP_TAG="${REGISTRY:+$REGISTRY/}$IMAGE_NAME:$(date +%Y%m%d-%H%M%S)"
        docker tag "$FULL_IMAGE_NAME" "$TIMESTAMP_TAG"
        echo "Tagged as: $TIMESTAMP_TAG"
    fi
fi

# Push if requested
if [ "$PUSH" = true ]; then
    echo "Pushing image to registry..."
    docker push "$FULL_IMAGE_NAME"
    
    # Push additional tags
    for tag in $(docker images --format "table {{.Repository}}:{{.Tag}}" | grep "^${REGISTRY:+$REGISTRY/}$IMAGE_NAME:" | grep -v latest); do
        docker push "$tag"
    done
    
    echo "Push completed successfully!"
fi

echo "Done!" 