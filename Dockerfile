FROM gcr.io/distroless/static-debian12

COPY ./gerry /
WORKDIR /app

CMD [ "/gerry", "start", "-c", "/app/config.yaml" ]
