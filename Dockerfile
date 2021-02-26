FROM golang:1.16 as build

WORKDIR /app

COPY . .

RUN make build

FROM ubuntu:20.04

WORKDIR /app

COPY --from=build /app/configs ./configs
COPY --from=build /app/api .

EXPOSE 8080

CMD [ "./api" ]