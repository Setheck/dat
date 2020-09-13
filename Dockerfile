FROM golang:alpine AS builder

RUN apk update && apk add build-base make git musl
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build
WORKDIR /app
RUN cp /build/dat ./dat

FROM alpine

RUN apk --no-cache --update add tzdata
COPY --chown=65534:0 --from=builder /app /app
USER 65534

WORKDIR /data
ENTRYPOINT ["/app/dat"]