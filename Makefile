build:
	go build -o bin/simser .

run: build
	go run bin/simser 

test:
	go test -v ./... -count=1

clean:
	rm -r ./bin