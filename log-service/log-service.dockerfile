FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY ./bin/log-app ./
COPY ./configs ./configs

CMD [ "/app/log-app" ]