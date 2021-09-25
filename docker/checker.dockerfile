FROM golang AS build-dev

RUN adduser --disabled-password -u 10000 checker
RUN mkdir /build/ && chown checker /build/
USER checker

WORKDIR /build/
ADD . /build/

RUN CGO_ENABLED=0 go build -o /build/checker ./cmd/checker

FROM alpine

RUN adduser -D -u 10000 checker
USER checker

WORKDIR /

COPY --from=build-dev /build/checker /checker

EXPOSE 9000

CMD [ "/checker" ]
