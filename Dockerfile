FROM golang:latest as builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monerohub .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/monerohub .

EXPOSE 8080

CMD ["./monerohub"]
