all:
	go build -o geoguessr

install: all
	cp geoguessr /usr/local/bin/

uninstall:
	rm /usr/local/bin/geoguessr

.PHONY: run
run: all
	./geoguessr

.PHONY: clean
clean:
	go clean
