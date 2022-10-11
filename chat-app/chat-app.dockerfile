FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY ./.bin/chat-app .
COPY ./configs ./configs

CMD [ "/app/chat-app" ]