#=============#
# TheQuestion #
# ============#

all: get install

fmt:
	@ go fmt ./...

get:
	@ govendor sync

install:
	@ go install ./...

test:
	@ go test -cover -v `go list ./... | grep -v /vendor/`

