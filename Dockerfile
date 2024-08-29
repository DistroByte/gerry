FROM gcr.io/distroless/static-debian12

COPY ./bin/gerry /
WORKDIR /app

CMD [ "/gerry", "start", "-c", "/app/config.yaml" ]
