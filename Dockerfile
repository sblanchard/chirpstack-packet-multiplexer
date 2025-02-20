FROM golang:1.13-alpine AS development

LABEL org.opencontainers.image.source="https://github.com/sblanchard/chirpstack-packet-multiplexer"
LABEL org.opencontainers.image.description="Forward Semtech packet-forwarder data to multiple servers."
LABEL org.opencontainers.image.authors="sblanchard"

ENV PROJECT_PATH=/chirpstack-packet-multiplexer
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM alpine:latest AS production

WORKDIR /root/
RUN apk --no-cache add tzdata
COPY --from=development /chirpstack-packet-multiplexer/build .
ENTRYPOINT ["./chirpstack-packet-multiplexer"]


