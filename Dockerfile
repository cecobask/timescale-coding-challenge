FROM golang:1.22-alpine as builder
WORKDIR /mnt/timescale-coding-challenge
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o build/ts main.go

FROM scratch
COPY --from=builder /mnt/timescale-coding-challenge/build/ts /usr/local/bin/ts
ENTRYPOINT ["/usr/local/bin/ts"]
CMD ["--help"]
