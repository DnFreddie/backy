FROM golang:1.22
WORKDIR /scripts

COPY script.sh .
RUN chmod +x /scripts/script.sh
RUN /scripts/script.sh

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main .


CMD sh -c './main dot && ls -a /root/.config && ls -a /tmp/test_data'

