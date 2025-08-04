.PHONY: all

all:
	rm -rfv dist
	go run main.go
	python3 -m http.server --bind localhost --directory ./dist 8080
	
