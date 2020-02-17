REPO    := 113379206287.dkr.ecr.us-east-1.amazonaws.com/
NAME    := dannylesnik/simple-json-server
TAG     := $$(git describe --tags `git rev-parse HEAD`)
IMG     := ${NAME}:${TAG}
LATEST  := ${NAME}:latest

all: build tag push

build:
		@docker build -t "$(IMG)" .

tag:
		@docker tag "$(IMG)" "$(REPO)${NAME}:${TAG}"

tag-latest:
		@docker tag "$(IMG)" "$(REPO)${LATEST}"

push:
		@docker push "$(REPO)${NAME}:${TAG}"

yaml-qa:
		@kustomize edit set image "$(REPO)${NAME}"="$(REPO)${NAME}:${TAG}"
		@kustomize edit add configmap http-json-server-config --from-file=deployment/qa/config.properties
		@kustomize build . -o deploy.yaml

yaml-prod:
		@kustomize edit set image "$(REPO)${NAME}"="$(REPO)${NAME}:${TAG}"
		@kustomize edit add configmap http-json-server-config --from-file=deployment/production/config.properties
		@kustomize build . -o deploy.yaml