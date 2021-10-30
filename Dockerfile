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

COPY --from=build /app/ ./

EXPOSE 8080

CMD ["./main"]


