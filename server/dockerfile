FROM golang:1.16-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache make

LABEL maintainer="Israel A."

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8081

CMD ["make", "build-run"]