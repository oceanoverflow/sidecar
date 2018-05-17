sidecar: 
	dep ensure
	go build -o sidecar .

clean:
	rm -f sidecar