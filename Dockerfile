FROM golang:1.16 as builder

WORKDIR /workspace
COPY ./ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o kamwiel main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kamwiel .
USER nonroot:nonroot

ENV PORT=3000 \
    GIN_MODE=release

ENTRYPOINT ["/kamwiel"]
