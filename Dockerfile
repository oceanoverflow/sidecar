FROM golang:alpine
COPY . /home
WORKDIR /home
RUN go build -o sidecar ./main.go
CMD ./start-agent.sh