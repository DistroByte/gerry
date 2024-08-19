FROM golang:1.23-bookworm AS builder

WORKDIR /go/src/gerry
COPY go.mod go.sum ./
COPY vendor ./
COPY . .

RUN make build

FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/src/gerry/build/gerry /
WORKDIR /app

CMD [ "/gerry", "start", "/app/config.yaml" ]
