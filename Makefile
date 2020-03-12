M = $(shell printf "\033[34;1m▶\033[0m")


.PHONY: setup-internal-dep
setup-internal-dep:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.6

.PHONY: schema
schema: $(info $(M) Generating GQL schema and resolvers)
	go run github.com/99designs/gqlgen generate

.PHONY: local-nginx
local-nginx: $(info $(M) Starting nginx locally)
	killall nginx
	nginx -c `pwd`/nginx/sites-available/local.delphishq.com