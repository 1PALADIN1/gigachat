FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY ./bin/auth-app .
COPY ./configs ./configs

CMD [ "/app/auth-app" ]