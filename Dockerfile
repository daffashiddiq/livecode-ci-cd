# STAGE 1
FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o livecode-golang

# STAGE 2
FROM alpine
WORKDIR /app
COPY --from=build /src/livecode-golang /app
ENTRYPOINT [ "/app/livecode-golang" ]