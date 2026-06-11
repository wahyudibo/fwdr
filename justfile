test:
    go test ./... -race -count=1

lint:
    golangci-lint run

build:
    CGO_ENABLED=0 go build -ldflags="-w -s" -o fwdr ./cmd/fwdr

docker-build tag="latest":
    DOCKER_BUILDKIT=1 docker build --tag fwdr:{{tag}} .

helm-install release="fwdr" namespace="default":
    helm upgrade --install {{release}} ./deployment/helm --namespace {{namespace}} --create-namespace

helm-uninstall release="fwdr" namespace="default":
    helm uninstall {{release}} --namespace {{namespace}}

helm-lint:
    helm lint ./deployment/helm

# Push docker image to a registry — e.g. just docker-push ghcr.io/youruser 1.0.0
docker-push registry tag="latest":
    docker tag fwdr:{{tag}} {{registry}}/fwdr:{{tag}}
    docker push {{registry}}/fwdr:{{tag}}

# Package helm chart into a .tgz — version must be semver without v prefix
helm-package version="0.1.0":
    helm package ./deployment/helm --version {{version}} --app-version {{version}}

# Push packaged helm chart to an OCI registry — e.g. just helm-push ghcr.io/youruser 0.1.0
helm-push registry version="0.1.0":
    helm push fwdr-{{version}}.tgz oci://{{registry}}
