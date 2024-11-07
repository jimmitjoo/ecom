.PHONY: validate-api
validate-api:
	swagger validate ./docs/swagger.yaml
	go test ./src/tools/validate/...

.PHONY: generate-client
generate-client:
	swagger generate client \
		-f docs/swagger.yaml \
		-c src/client/swagger.config.yaml \
		-t src/client \
		--skip-validation

.PHONY: clean-client
clean-client:
	rm -rf src/client/operations src/client/models