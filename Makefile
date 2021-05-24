shell := /bin/bash

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: build
build:
	cd cmd/api && go build -o indahaus

.PHONY: migrate
migrate:
	go run cmd/admin/main.go migrate

.PHONY: resetdb
resetdb:
	rm indahaus.db
	go run cmd/admin/main.go migrate

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	# see https://github.com/99designs/gqlgen/issues/1483 and https://github.com/golang/go/issues/44129
	go get github.com/99designs/gqlgen/cmd@v0.13.0 && go run github.com/99designs/gqlgen generate

# ===================================================================
# docker build
.PHONY: docker-build
docker-build:
	docker build --progress=plain \
		-f config/Dockerfile \
		-t indahaus-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%d%H:%M:%SZ"` \
		--build-arg PORT=8080 \
		.

.PHONY: docker-run
docker-run:
	docker run --rm -p 8080:8080 -e AUTH_PASSWORD=supersecret -e AUTH_USERNAME=secureworks indahaus-amd64:1.0

# ===================================================================
# k8s dev
.PHONY: kind-up kind-down kind-load api-up postgres-up helm-up helm-down up down cluster-info update-api
kind-up:
	-kind create cluster --image kindest/node:v1.20.2 --name indahaus
	kubectl cluster-info --context kind-indahaus

kind-down:
	kind delete cluster --name indahaus

kind-load:
	kind load docker-image indahaus-amd64:1.0 --name indahaus

helm-up: 
	helm install api config/charts

helm-down: 
	helm delete api

# bring up the full application for dev
up: docker-build kind-up kind-load helm-up

# tear down the full application including the kind cluster
down: helm-down kind-down

# view info about deployed helm charts and check that pods are up and running
cluster-info:
	helm ls
	kubectl get pods --watch

# rebuild the image, load it into kind and delete the old pod so k8s will restart with code updates
update-api: docker-build
	kind load docker-image indahaus-amd64:1.0 --name indahaus
	kubectl delete pods -lapp.kubernetes.io/name=indahaus

.PHONY: logs
logs:
	kubectl logs -l app.kubernetes.io/name=indahaus -f

# forward the api port
.PHONY: port-forward-api
port-forward-api:
	kubectl get pods -l app.kubernetes.io/name=indahaus | grep api | cut -d' ' -f1 | xargs -I{} kubectl port-forward {} 8080:8080

# forward the debug/metrics port
.PHONY: port-forward-debug
port-forward-debug:
	kubectl get pods -l app.kubernetes.io/name=indahaus | grep api | cut -d' ' -f1 | xargs -I{} kubectl port-forward {} 4000:4000