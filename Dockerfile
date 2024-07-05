# Use an appropriate Go version tag
FROM golang:1.22.2-alpine
ENV USER=test
RUN set -eux; \
    apk update && \
    apk add --no-cache git

RUN adduser -D -h /home/${USER} ${USER}



COPY . /app

WORKDIR /app


RUN go mod tidy

USER ${USER}


