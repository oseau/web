ARG IMAGE_NODE
FROM ${IMAGE_NODE}

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
RUN pnpm config set store-dir "$PNPM_HOME/.pnpm-store"

WORKDIR /usr/src/app
