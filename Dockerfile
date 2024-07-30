FROM golang:latest AS build

WORKDIR /build

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

ENV HOME=/app
COPY --from=build /build/dnschecker /bin/

ENTRYPOINT ["dnschecker"]