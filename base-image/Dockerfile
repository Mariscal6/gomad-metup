# Build the simple-controller binary
FROM golang:1.23 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY go.mod go.mod
RUN go mod download
RUN go generate ./...
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o gopher-pprof main.go
FROM scratch
WORKDIR /
COPY --from=builder /workspace/gopher-pprof .
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/gopher-pprof"]
