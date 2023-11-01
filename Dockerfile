FROM golang:1.20

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN apt update && apt upgrade -y

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air