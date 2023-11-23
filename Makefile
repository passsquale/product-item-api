.PHONY: run
run:
	go run cmd/product-item-api/main.go

.PHONY: test
test:
	go test internal/app/retranslator/retranslator_test.go