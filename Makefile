.PHONY: all test secret image query deps run service-new


test:
	go test ./...

test-slack:
	go test --run TestSlackHandler ./pkg/handler/slack/ -v

test-rest:
	go test --run TestRestHandler ./pkg/handler/rest/ -v

test-chat:
	go test --run TestChatHandler ./pkg/handler/chat/ -v

run:
	go run ./cmd/*.go

secret:
	kubectl delete secret twfeel-secrets -n demo --ignore-not-found=true
	kubectl create secret generic twfeel-secrets -n demo \
		--from-literal=ACCESS_TOKENS="rest:${TW_FEEL_REST_TOKEN};chat:${TW_FEEL_CHAT_TOKEN};slack:${TW_FEEL_SLACK_TOKEN}" \
		--from-literal=REDIS_PASS=$(REDIS_PASS) \
		--from-literal=T_CONSUMER_KEY=$(TK_CONSUMER_KEY) \
		--from-literal=T_CONSUMER_SECRET=$(TK_CONSUMER_SECRET) \
		--from-literal=T_ACCESS_TOKEN=$(TK_ACCESS_TOKEN) \
		--from-literal=T_ACCESS_SECRET=$(TK_ACCESS_SECRET)

image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/twfeel:latest

docker:
	docker build -t twfeel .

service:
	kubectl apply -f config/service.yaml

service-new:
	kubectl delete -f config/service.yaml --ignore-not-found=true
	kubectl apply -f config/service.yaml

query:
	curl -H "Content-Type: application/json" -X GET \
		"https://twfeel.demo.knative.tech/v1/feel/knative?token=${TW_FEEL_TOKEN}" \
		| jq "."

redis-connect:
	kubectl exec -it redis-0 -n demo /bin/sh

deps:
	go get -u cloud.google.com/go
	go get -u github.com/dghubble/go-twitter
	go get -u github.com/dghubble/oauth1
	go get -u github.com/gin-gonic/gin
	go get -u github.com/go-redis/cache
	go get -u github.com/go-redis/redis
	go get -u github.com/google/uuid
	go get -u github.com/stretchr/testify
	go get -u github.com/vmihailenco/msgpack
	go mod tidy