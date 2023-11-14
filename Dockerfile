FROM golang:1.20-alpine3.18 AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o scaler -ldflags '-w -extldflags "-static"' .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/scaler /usr/local/bin/scaler

ENTRYPOINT ["scaler"]