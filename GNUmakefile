VERSION=1.0.0
LOCAL_PROVIDER_PATH="$$HOME/.terraform.d/plugins/registry.terraform.io/kayteh/podio/${VERSION}/$$(go env GOOS)_$$(go env GOARCH)"

default: build

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: build
build:
	mkdir -p ${LOCAL_PROVIDER_PATH}
	go build -v -o ${LOCAL_PROVIDER_PATH}/terraform-provider-podio_v${VERSION}
	@echo "Built podio provider v${VERSION} for $$(go env GOOS)_$$(go env GOARCH) at ${LOCAL_PROVIDER_PATH}/terraform-provider-podio_v${VERSION}"
