# syntax=docker/dockerfile:1

FROM golang:1.19.0-alpine3.16 as builder

RUN apk add --no-cache git

RUN mkdir -p /go/src/github.com/syoder89

WORKDIR /go/src/github.com/syoder89

RUN git clone https://github.com/syoder89/tank-monitor

WORKDIR /go/src/github.com/syoder89/tank-monitor

RUN go build -o /tank-monitor

FROM golang:1.19.0-alpine3.16
COPY --from=builder /tank-monitor /tank-monitor
CMD [ "/tank-monitor" ]
