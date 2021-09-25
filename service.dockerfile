FROM golang AS build-dev

RUN adduser --disabled-password -u 10000 checker
RUN mkdir /build_checker/ && chown checker /build_checker/
USER checker

WORKDIR /build_checker/
ADD . /build_checker/

RUN CGO_ENABLED=0 go build -o /build_checker/checker ./cmd

FROM alpine

RUN adduser -D -u 10000 checker
USER checker

WORKDIR /

COPY --from=build-dev /build_checker/checker /checker

EXPOSE 9000

CMD [ "/checker" ]
