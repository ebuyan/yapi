#--------------------------------------
FROM golang:1.18.1-alpine3.14 as build

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /app/yapi /app/cmd/yapi.go

#--------------------------------------
FROM alpine:3.14 as app

WORKDIR /app

COPY --from=build /app/yapi /app/yapi
COPY --from=build /app/.env.local /app/.env.local

RUN chmod +x /app/yapi

EXPOSE 8001

ENTRYPOINT /app/yapi
