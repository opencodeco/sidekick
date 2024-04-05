sidekick: clean
	go build -o ./bin/sidekick ./cmd/sidekick

clean:
	@if [ -f ./bin/sidekick ]; then rm ./bin/sidekick; fi
