FROM golang:alpine3.19 AS build
WORKDIR /app
COPY ./go.mod /app/ 
COPY ./go.sum /app/
COPY ./*.go /app/
RUN go build -o lister .

FROM alpine:latest as security_provider
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot


FROM scratch
COPY --from=security_provider /etc/passwd /etc/passwd
USER nonroot
COPY --from=build /app/lister /app/lister
WORKDIR /app
ENTRYPOINT [ "./lister" ]