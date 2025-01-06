# syntax=docker/dockerfile:1

FROM alpine:latest

# These are automatically set by the build command
ARG TARGETOS
ARG TARGETARCH

COPY ./artifacts/ ./artifacts/
COPY VERSION.txt ./

RUN <<EOF
ARCH=$(case "${TARGETARCH}" in
    amd64) echo "x86_64" ;;
    arm64) echo "aarch64" ;;
    *) echo "Unsupported architecture: ${TARGETARCH}" >&2; exit 1 ;;
esac)
VERSION=$(cat VERSION.txt)
mv ./artifacts/dqlite-vip-${VERSION}-${TARGETOS}-${ARCH} /dqlite-vip
EOF

FROM scratch
COPY --chmod=775 --from=0 /dqlite-vip /dqlite-vip
ENTRYPOINT ["/dqlite-vip"]