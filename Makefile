all:
	go build -o geoguessr

.PHONY: run
run: all
	./geoguessr

.PHONY: clean
clean:
	go clean
