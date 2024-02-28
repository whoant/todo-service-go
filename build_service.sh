IMAGE_NAME=todo-service
CACHED_BUILD=$1

if [[ -n "$CACHED_BUILD" ]]; then
  echo "Docker building cached image..."
  docker rmi ${IMAGE_NAME}-cached ${IMAGE_NAME}
  docker build -t ${IMAGE_NAME}-cached -f cache.Dockerfile .
fi

echo "Docker building main image..."
docker build -t todo-service:latest -f Dockerfile .

echo "Done!!"
