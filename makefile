all: richman

.PHONY: richman
richman:
	go build -o richman

clean:
	rm richman