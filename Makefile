GITHASH:=$(shell git rev-parse --short HEAD)
BUILDTIME:=$(shell TZ=UTC date +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS='-w -extldflags "-static" -s -X github.com/scm-manager/cli/pkg.commitHash=${GITHASH} -X github.com/scm-manager/cli/pkg.buildTime=${BUILDTIME}'

.PHONY:
build:
	go build -a -tags netgo -ldflags ${LDFLAGS} -o scm scm.go
