sidekick:
	@go build -o ./bin/sidekick ./cmd/sidekick

.PHONY: install
install:
	@go install ./cmd/sidekick
