FROM golang:1.22.2-alpine

ENV USER=test

RUN set -eux; \
    apk update && \
    apk add --no-cache git && \
    adduser -D -h /home/${USER} ${USER}

COPY . /app

WORKDIR /app

RUN go mod download

RUN go mod tidy

RUN chown -R ${USER}:${USER} /app

USER ${USER}

# Optionally, specify a command to run your application or tests
# CMD ["go", "run", "main.go"]  # Uncomment and modify as needed

