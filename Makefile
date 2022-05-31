GO := go

all: build build-plugins

build:
	$(GO) build \
		-v \
		-o server \
		.

build-plugins:
	mkdir -p plugins
	$(GO) build \
		-v \
		-o plugins/ \
		-buildmode=plugin \
		./plugin_*
