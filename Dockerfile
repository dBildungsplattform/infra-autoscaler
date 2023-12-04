FROM golang:1.20-alpine3.18 AS build_deps

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o infra-autoscaler -ldflags '-w -extldflags "-static"' .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/infra-autoscaler /usr/local/bin/infra-autoscaler
COPY config/scaler_config.yml config/scaler_config.yml
ENTRYPOINT ["infra-autoscaler"]