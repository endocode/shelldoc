all:	VERSION build

VERSION:
	git describe  --always --tags --abbrev=7 HEAD > VERSION

VersionString = $(shell head -n 1 VERSION | tr -d '\n')

build				: VERSION
	cd cmd/shelldoc2 && go build -ldflags '-X github.com/endocode/shelldoc/pkg/version.versionString=${VersionString}'
