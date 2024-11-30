.PHONY: test
test:
	@go test ./... -tags=unit

.PHONY: testjson
testjson:
	@go test ./... -tags=unit -json > testresults.json
