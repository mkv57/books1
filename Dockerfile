FROM golang:1.23.4-alpine3.20 as builder

RUN apk update && apk add --no-cache ca-certificates 
# RUN apk update && apk add --no-cache ca-certificates go-task
WORKDIR /my_proekt


COPY go.mod ./

COPY go.sum ./

RUN go mod download


COPY . .

# docker build -t books1 .
# docker compose up --build

RUN go build ./cmd/libs

FROM scratch

WORKDIR /bin

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /etc/task/go-task /etc/task/go-task // если надо
COPY --from=builder /my_proekt/libs /libs

#COPY --from=builder /my_proekt/config.yml ./config.yml
COPY --from=builder /my_proekt/app.log ./app.log
COPY --from=builder /my_proekt/migrate ./migrate

ENTRYPOINT ["/libs"]  

# docker build -t registry.academy.the-guild.tech/book_mkv:0.1.0 .
# docker push registry.academy.the-guild.tech/book_mkv:0.1.0 
# docker build -t <container_name>:<tag> . (с точкой в конце)
# docker push <container_name>:<tag>        (без точки)