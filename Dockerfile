ARG GO_VERSION=1

FROM golang:${GO_VERSION}-alpine as builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app .

FROM alpine
# add required tzdata. ref: https://stackoverflow.com/a/62159987/2496217
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=builder /app /app/
COPY static/ /app/static/
CMD ["/app/app"]
