FROM golang:1.10-alpine3.7 as builder

# Set GOPATH/GOROOT environment variables
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# go get all of the dependencies
RUN apk update \
  && apk add --virtual build-deps musl-dev curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Set up app
WORKDIR /go/src/gostartup
COPY . .
RUN dep ensure
RUN go build -o api src/api.go
RUN go build -o callback src/callback.go
RUN go build -o sync src/sync.go

FROM alpine:3.7

COPY --from=builder /go/src/gostartup/api /
COPY --from=builder /go/src/gostartup/callback /
COPY --from=builder /go/src/gostartup/sync /
COPY --from=builder /go/src/gostartup/.env* /
COPY /wait-for /

EXPOSE 3000
EXPOSE 5000

ENTRYPOINT ["/wait-for", "db:5432", "--"]

CMD ["/api"]
