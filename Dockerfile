FROM golang:1.22-alpine as build
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o gerry .

FROM alpine:3
RUN addgroup -g 1000 gerry && adduser -u 1000 -G gerry -D gerry
USER gerry
COPY --from=build /src/gerry /bin
COPY --from=build /src/plugins /data/plugins
VOLUME [ "/data/gerry" ]
ENV GERRY_PLUGINS_PATH=/data/plugins
ENV GERRY_WATCHER_ENABLED=0
CMD [ "/bin/gerry" ]