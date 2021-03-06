FROM golang:1.16 as build

WORKDIR /go/src

COPY . .

RUN make build

FROM debian:buster-slim

WORKDIR /app

COPY --from=build /go/src/configs ./configs
COPY --from=build /go/src/api .

EXPOSE 8080

CMD [ "./api" ]