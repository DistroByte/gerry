# FROM golang:1.22-alpine as build
# WORKDIR /src
# COPY go.mod .
# COPY go.sum .
# RUN go mod download
# COPY . .
# RUN 

FROM alpine:3
RUN addgroup -g 1000 gerry && adduser -u 1000 -G gerry -D gerry
USER gerry
COPY ./build/gerry /bin
CMD [ "/bin/gerry" ]