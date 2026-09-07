FROM golang:1.27.1-alpine3.24 AS build_deps

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o infra-autoscaler -ldflags '-w -extldflags "-static"' .

FROM alpine:3.24

RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/infra-autoscaler /usr/local/bin/infra-autoscaler
COPY config/scaler_config.yml config/scaler_config.yml

RUN chown -R nobody:nogroup /usr/local/bin/infra-autoscaler config/scaler_config.yml
USER nobody

ENTRYPOINT ["infra-autoscaler"]
