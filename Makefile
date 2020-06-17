M = $(shell printf "\033[34;1m▶\033[0m")
NOW = ${shell date +%s}


.PHONY: setup-internal-dep
setup-internal-dep:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.0

.PHONY: schema
schema: $(info $(M) Generating GQL schema and resolvers)
	go run github.com/99designs/gqlgen generate

.PHONY: local-nginx
local-nginx: $(info $(M) Starting nginx locally)
	killall nginx || true > /dev/null
	nginx -c `pwd`/nginx/sites-available/local.delphishq.com

run-local:
	DELPHIS_ENV=local go run server.go

run-local-use-aws:
	AWS_PROFILE=delphis DELPHIS_ENV=local_use_aws go run server.go

build:
	go build -o delphis_server

build-and-deploy-docker: get-ecr-creds
	docker build -t delphisbe .
	docker tag delphisbe:latest 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:latest
	docker tag delphisbe:latest 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:${NOW}
	docker push 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:latest
	docker push 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:${NOW}

update-service:
	aws ecs update-service --cluster delphis-cluster --service delphis-service --force-new-deployment --profile delphis

get-ecr-creds:
	aws ecr --profile delphis get-login-password --region us-west-2 | docker login --username AWS --password-stdin 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe

start-db:
	pg_ctl -D '/usr/local/var/postgresql@11/data' start

.PHONY: mocks
mocks:
	${GOPATH}/bin/mockery -output ./mocks -name Datastore -dir ./internal/datastore -case underscore
