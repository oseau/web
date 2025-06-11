ARG IMAGE_GO
FROM ${IMAGE_GO}
# https://docs.docker.com/build/building/variables/#scoping
# Consume the build argument in the build stage
ARG VERSION_STATICCHECK
ARG VERSION_REVIVE

RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends entr
rm -rf /var/lib/apt/lists/*
go install honnef.co/go/tools/cmd/staticcheck@${VERSION_STATICCHECK}
go install github.com/mgechev/revive@${VERSION_REVIVE}
EOF

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download
