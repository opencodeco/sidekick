sidekick: clean
	go build github.com/opencodeco/sidekick

clean:
	if [ -f sidekick ]; then rm sidekick; fi
