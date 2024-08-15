FROM golang:1.22.2-alpine

ENV USER=test
ENV GO111MODULE=on

RUN set -eux; \
    apk update && \
    apk add --no-cache git tmux bash && \
    adduser -D -h /home/${USER} ${USER}

WORKDIR /home/${USER}/backy

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p /home/${USER}/.config/\
    && touch /home/${USER}/.config/LICENSE

RUN chown -R ${USER}:${USER} /home/${USER} 
#&& \
#chmod -R 777 /home/${USER}

USER ${USER}

