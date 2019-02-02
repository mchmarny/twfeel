
setup:
	export ACCESS_TOKEN=$(TW_FEEL_TOKEN)
	export T_CONSUMER_KEY=$(T_CONSUMER_KEY)
	export T_CONSUMER_SECRET=$(T_CONSUMER_SECRET)
	export T_ACCESS_TOKEN=$(T_ACCESS_TOKEN)
	export T_ACCESS_SECRET=$(T_ACCESS_SECRET)

secret:
	kubectl delete secret twfeel-secrets -n demo --ignore-not-found=true
	kubectl create secret generic twfeel-secrets -n demo \
		--from-literal=ACCESS_TOKEN=$(TW_FEEL_TOKEN) \
		--from-literal=T_CONSUMER_KEY=$(TK_CONSUMER_KEY) \
		--from-literal=T_CONSUMER_SECRET=$(TK_CONSUMER_SECRET) \
		--from-literal=T_ACCESS_TOKEN=$(TK_ACCESS_TOKEN) \
		--from-literal=T_ACCESS_SECRET=$(TK_ACCESS_SECRET)

image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/twfeel:latest

service:
	kubectl apply -f service.yaml

query:
	curl -H "Content-Type: application/json" -X GET \
		"https://twfeel.demo.knative.tech/v1/feel/knative?token=${TW_FEEL_TOKEN}" \
		| jq "."

