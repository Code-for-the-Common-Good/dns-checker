VERSION ?= $(shell git describe --tags --always)
PACKAGENAME ?= dnschecker
LDFLAGS ?= "-w -s"

docker: 
	docker-compose build
	docker-compose up

build:  ## Build git-chglog
	$(MAKE) --no-print-directory log-$@
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o $(PACKAGENAME)

changelog: build   ## Generate changelog
	@ $(MAKE) --no-print-directory log-$@
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md

release: changelog   ## Release a new tag
	@ $(MAKE) --no-print-directory log-$@
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
	git tag $(VERSION)
	git push origin main $(VERSION)

help:   ## Display this help
	@awk \
		-v "col=\033[36m" -v "nocol=\033[0m" \
		' \
			BEGIN { \
				FS = ":.*##" ; \
	printf "Usage:\n  make %s<target>%s\n", col, nocol \
      } \
      /^[a-zA-Z_-]+:.*?##/ { \
        printf "  %s%-12s%s %s\n", col, $$1, nocol, $$2 \
      } \
      /^##@/ { \
        printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol \
      } \
    ' $(MAKEFILE_LIST)

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'