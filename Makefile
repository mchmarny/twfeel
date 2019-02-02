.PHONY: all test secret image query update run


test:
	go test ./... -v

run:
	go run cmd/*.go

secret:
	kubectl delete secret twfeel-secrets -n demo --ignore-not-found=true
	kubectl create secret generic twfeel-secrets -n demo \
		--from-literal=ACCESS_TOKEN=$(TW_FEEL_TOKEN) \
		--from-literal=REDIS_PASS=$(REDIS_PASS) \
		--from-literal=T_CONSUMER_KEY=$(TK_CONSUMER_KEY) \
		--from-literal=T_CONSUMER_SECRET=$(TK_CONSUMER_SECRET) \
		--from-literal=T_ACCESS_TOKEN=$(TK_ACCESS_TOKEN) \
		--from-literal=T_ACCESS_SECRET=$(TK_ACCESS_SECRET)

image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/twfeel:latest

service:
	kubectl apply -f config/service.yaml

query:
	curl -H "Content-Type: application/json" -X GET \
		"https://twfeel.demo.knative.tech/v1/feel/knative?token=${TOKEN}" \
		| jq "."

update:
	go get -u github.com/gin-gonic/gin
	go get -u cloud.google.com/go
	go get -u github.com/dghubble/go-twitter
	go get -u github.com/dghubble/oauth1
	go get -u github.com/gin-gonic/gin
	go get -u github.com/go-redis/redis
	go get -u github.com/go-redis/cache
	go get -u github.com/google/uuid
	go get -u github.com/stretchr/testify
	go mod tidy
	go mod vendor