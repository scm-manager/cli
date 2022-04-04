VERSION:=0.1.0
GIHASH:=$(shell git rev-parse --short HEAD)
BUILDTIME:=$(shell TZ=UTC date +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS='-w -extldflags "-static" -s -X github.com/scm-manager/cli/pkg.version=${VERSION} -X github.com/scm-manager/cli/pkg.gitHash=${GIHASH} -X github.com/scm-manager/cli/pkg.buildTime=${BUILDTIME}'

.PHONY:
build:
	go build -a -tags netgo -ldflags ${LDFLAGS} -o scm scm.go
