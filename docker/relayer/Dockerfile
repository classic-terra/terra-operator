FROM golang:1.19-alpine AS builder

WORKDIR /src/

ARG COMMIT_HASH=v2.2.0

RUN apk add --update git make gcc linux-headers libc-dev eudev-dev
RUN git clone https://github.com/cosmos/relayer.git \
    && cd relayer \
    && git checkout $COMMIT_HASH \
    && make install

FROM public.ecr.aws/p5q2r9h7/node:alpine3.17
COPY --from=builder /go/bin/rly /usr/local/bin/
RUN apk add bash vim \
    && addgroup -g 1000 relayer \
    && adduser -S -h /home/relayer -D relayer -u 1000 -G relayer

USER 1000
WORKDIR /home/relayer

# Create startup script that takes a path as a parameter
RUN echo 'rly start $1 -p events' > start.sh

COPY ./chain-config-format.json /chain-config-format.json

CMD ["/entrypoint.sh"]