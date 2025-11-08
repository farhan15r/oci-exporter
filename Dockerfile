FROM docker.io/library/golang:1.25.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o oci-exporter ./main.go


FROM docker.io/library/alpine:3.20

WORKDIR /app

COPY --from=build /app/oci-exporter .

CMD ["./oci-exporter"]