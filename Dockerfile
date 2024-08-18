FROM golang:1.23-bookworm AS builder

WORKDIR /go/src/gerry
COPY go.mod go.sum ./
COPY vendor ./
COPY . .

RUN make build

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /go/src/gerry/build/gerry /

CMD [ "/gerry", "start", "/etc/gerry/config.yaml" ]
