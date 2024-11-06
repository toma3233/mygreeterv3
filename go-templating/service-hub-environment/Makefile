REGISTRY_NAME = servicehubregistry
IMG_NAME = service_hub_environment
DATE = $(shell date +%Y%m%d)

image:
	if ! docker buildx ls | grep -q svchubbuilder; then \
		docker buildx create --name svchubbuilder --driver docker-container --bootstrap --use; \
	else \
		docker buildx use svchubbuilder; \
	fi
	az acr login --name $(REGISTRY_NAME)
	docker buildx build --platform linux/amd64,linux/arm64 --tag $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(DATE) --push -f Dockerfile ..
	docker buildx use --builder default
	docker buildx rm --builder svchubbuilder

single-arch-image:
	docker build -t $(IMG_NAME) -f Dockerfile ..
	docker tag $(IMG_NAME) $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(DATE)
	az acr login --name $(REGISTRY_NAME)
	docker push $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(DATE)