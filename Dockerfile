FROM golang:1.13.6-alpine3.10 AS builder

RUN apk add --no-cache \
        make

COPY . /tmp/build
WORKDIR /tmp/build

RUN make build

FROM alpine:3.10
COPY --from=builder /tmp/build/drlm-agent /
COPY agent.toml /

RUN mkdir /root/.bin
ENV PATH="/root/.bin:${PATH}"

CMD [ "/drlm-agent", "-j" ]