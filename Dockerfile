# build stage
FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/vinayr/go-garden
COPY . .
RUN apk add --update make && make

# final stage
FROM alpine
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/vinayr/go-garden/garden /app/
COPY --from=build-env /go/src/github.com/vinayr/go-garden/entrypoint.sh /app/
CMD ["./entrypoint.sh"]
