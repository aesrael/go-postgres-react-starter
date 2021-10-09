FROM node:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash make yarn

LABEL maintainer="Israel A."

WORKDIR /app

COPY package.json ./
COPY yarn.lock ./

RUN yarn

COPY . .

EXPOSE 3000

CMD ["yarn", "start"]