FROM golang:alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /ratingd ./cmd/ratingd

FROM alpine:3.20
COPY --from=build /ratingd /ratingd
ENTRYPOINT ["/ratingd"]
