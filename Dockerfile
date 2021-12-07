##
## Build
##
FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

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
COPY --from=build /app/serviceConfig.json ./
COPY --from=build /app/wait-for-it.sh ./

RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash

RUN chmod +x ./wait-for-it.sh

EXPOSE 6969

CMD ["./wait-for-it.sh", "db:3306", "--", "./main", "--config", "serviceConfig.json"]


