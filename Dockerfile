##
## Build
##
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy
RUN apk add build-base
RUN apk add libwebp-dev

COPY . .

RUN go build -o main main.go

##
## RUN
##
FROM alpine

WORKDIR /app

ARG REVISION_ID

LABEL revision_id=${REVISION_ID}

COPY --from=build /app/main ./
COPY --from=build /app/wait-for-it.sh ./

RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash

RUN apk add libwebp-dev

RUN chmod +x ./wait-for-it.sh

EXPOSE 6969

CMD ["./wait-for-it.sh", "db:3306", "--", "./main", "--config", "serviceConfig.json"]


