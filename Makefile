SHELL := /bin/bash

SERVICE_NAME := quasar-mcp-server
REGION := asia-northeast1
PROJECT_ID := $(shell gcloud config get-value project)
TARGET_URL := $(shell gcloud run services describe quasar --region ${REGION} --format 'value(status.url)' --project ${PROJECT_ID})
IMAGE := ${REGION}-docker.pkg.dev/${PROJECT_ID}/${SERVICE_NAME}/app
TAG := latest

update:
	go get -u ./...
	go mod tidy

install:
	go install

proxy:
	gcloud run services proxy quasar-mcp-server --region ${REGION} --port=3000

artifact:
	gcloud artifacts repositories create ${SERVICE_NAME} --repository-format=docker --REGION=${REGION} --project=${PROJECT_ID}

cloudbuild:
	gcloud builds submit --config cloudbuild.yaml --substitutions=_IMAGE=${IMAGE},_TAG=${TAG} .

build:
	gcloud auth configure-docker ${REGION}-docker.pkg.dev --quiet
	gcloud artifacts repositories list
	docker buildx build --platform=linux/amd64 -t ${IMAGE} .
	docker push ${IMAGE}

deploy:
	gcloud artifacts docker images describe ${IMAGE}
	gcloud run deploy --region ${REGION} --project ${PROJECT_ID} --image ${IMAGE} --set-env-vars=PROJECT_ID=${PROJECT_ID},TARGET_URL=${TARGET_URL} ${SERVICE_NAME}
	gcloud run services update-traffic ${SERVICE_NAME} --to-latest --region ${REGION} --project ${PROJECT_ID}
