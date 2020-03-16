FROM golang:1.14

MAINTAINER Brandon WELSCH <dev@brandon-welsch.eu>

# Install Reflex for Service Live Reload
RUN go get github.com/cespare/reflex

WORKDIR /app

COPY . .

RUN go build -mod vendor

CMD ["./game-servers"]
