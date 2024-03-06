# Makefile

.PHONY: build watch-css air dev

build:
	npm run build:css && go build -o bin/main cmd/main.go

watch-css:
	npm run watch:css

air:
	air

dev:
	$(MAKE) -j2 watch-css air