FROM golang:latest

RUN mkdir /app
COPY ./ /app
WORKDIR /app

RUN go mod download
RUN go build -o 'gigachat' /app/cmd/gigachat/main.go

CMD [ "/app/gigachat" ]