FROM golang:1.15-alpine as build

LABEL maintainer="Dušan Simić <dusan.simic1810@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server/server /app/server
COPY --from=build /app/db /app/db

ENV PORT=3000

EXPOSE ${PORT}

CMD [ "/app/server" ]
