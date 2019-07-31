all:	VERSION selftest

VERSION:
	git describe  --always --tags --abbrev=7 HEAD > VERSION

VersionString = $(shell head -n 1 VERSION | tr -d '\n')

build				: VERSION
	cd cmd/shelldoc && go build -ldflags '-X github.com/endocode/shelldoc/pkg/version.versionString=${VersionString}'

test:
	go test ./...

selftest			: build
	@echo "Running self-test of README.md and evaluating XML output with xmllint..." && \
		./cmd/shelldoc/shelldoc run -x results.xml README.md && \
		xmllint --noout --schema pkg/junitxml/jenkins-junit.xsd results.xml
