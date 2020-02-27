FROM golang:1.12-alpine3.7 as builder

# Set GOPATH/GOROOT environment variables
RUN mkdir -p /{{ cookiecutter.project_slug }}
ENV GOPATH /{{ cookiecutter.project_slug }}
ENV PATH $GOPATH/bin:$PATH

# go get all of the dependencies
RUN apk update \
  && apk add --virtual build-deps musl-dev curl git

# Set up app
WORKDIR /{{ cookiecutter.project_slug }}/
COPY . .
RUN go build -o {{ cookiecutter.project_slug }} src/{{ cookiecutter.project_slug }}.go

FROM alpine:3.7

COPY --from=builder /go/src/{{ cookiecutter.project_slug }}/api /
COPY --from=builder /go/src/{{ cookiecutter.project_slug }}/.env* /
COPY /wait-for /

EXPOSE 5000

CMD ["/{{ cookiecutter.project_slug }}"]
