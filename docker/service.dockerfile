FROM golang AS build-dev

RUN adduser --disabled-password -u 10000 service
RUN mkdir /build/ && chown service /build/
USER service

WORKDIR /build/
ADD . /build/

RUN CGO_ENABLED=0 go build -o /build/service ./cmd/service

FROM alpine

RUN adduser -D -u 10000 service
USER service

WORKDIR /

COPY --from=build-dev /build/service /service

EXPOSE 9000

CMD [ "/service" ]
