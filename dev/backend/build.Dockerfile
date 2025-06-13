ARG IMAGE_GO
ARG IMAGE_GO_PROD
ARG VERSION_XX
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${VERSION_XX} AS xx
FROM --platform=$BUILDPLATFORM ${IMAGE_GO} AS builder

COPY --from=xx / /
ARG TARGETPLATFORM
RUN xx-apt install -y libc6-dev gcc
ENV CGO_ENABLED=1

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG LDFLAGS
# Automatic platform ARG variables produced by Docker
ARG TARGETOS TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg <<EOF
xx-go --wrap
go build -ldflags "${LDFLAGS}" -a -o /backend ./cmd/backend/*.go
xx-verify /backend
EOF

FROM ${IMAGE_GO_PROD} AS prod
COPY --from=builder /backend /

CMD ["/backend"]
