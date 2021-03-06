M = $(shell printf "\033[34;1m▶\033[0m")
NOW = ${shell date +%s}
VER=$(shell git log -1 --pretty=format:"%H")

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
	docker tag delphisbe:latest 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:local-${VER}
	docker push 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:local-${VER}

create-task-def: build-and-deploy-docker
	rm -f task-curbuild-def.json
	VER=local-${VER} envsubst < task-definition.json > task-curbuild-def.json
	aws ecs register-task-definition --cli-input-json file://./task-curbuild-def.json --profile delphis --region us-west-2
	rm task-curbuild-def.json

update-service:
	aws ecs update-service --region us-west-2 --cluster delphis-cluster --service delphis-service --task-definition delphis-app-task --profile delphis

get-ecr-creds:
	aws ecr --profile delphis get-login-password --region us-west-2 | docker login --username AWS --password-stdin 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe

start-db:
	pg_ctl -D '/usr/local/var/postgresql@11/data' start

.PHONY: mocks
mocks:
	${GOPATH}/bin/mockery -output ./mocks -name Datastore -dir ./internal/datastore -case underscore
	${GOPATH}/bin/mockery -output ./mocks -name DelphisAuth -dir ./internal/auth -case underscore
	${GOPATH}/bin/mockery -output ./mocks -name TwitterClient -dir ./internal/twitter -case underscore
	${GOPATH}/bin/mockery -output ./mocks -name TwitterBackend -dir ./internal/twitter -case underscore

plan:
	@test "${env}" || (echo 'please pass in $$env' && exit)
	terraform fmt ./terraform
	@cd terraform && terraform init -backend=true -backend-config=backend/$(env).tfvars -get=true
	@cd terraform &&terraform plan -refresh=true -out $(env).plan
.PHONY: plan

apply:
	@test "${env}" || (echo 'please pass in $$env' && exit)
	@cd terraform && terraform apply -refresh $(env).plan
.PHONY: apply
