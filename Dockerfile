FROM golang:latest

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get github.com/go-chi/chi && \
    go get github.com/go-chi/render & \
    go get github.com/satori/go.uuid && \
    go get cloud.google.com/go/storage

RUN go build -o main
CMD ["/app/main"]