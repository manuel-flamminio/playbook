FROM golang:1.23 AS build

WORKDIR /go/src/app
COPY --link  . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
COPY ./seeder/seeder_file.json /data/seeder_file.json

CMD ["/app"]
