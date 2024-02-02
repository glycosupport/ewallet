FROM golang:1.21-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ewallet ./cmd/ewallet

FROM alpine:3.19

ARG USER=nonroot
ENV HOME /home/$USER

RUN apk add --update sudo

RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME

COPY --from=builder /app/ewallet ${HOME}/app/ewallet

ENTRYPOINT [ "./app/ewallet" ]