ARG IMAGE_NODE
FROM ${IMAGE_NODE}

ARG REPO_NAME
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV REPO_NAME=${REPO_NAME}
RUN corepack enable
RUN pnpm config set store-dir "$PNPM_HOME/.pnpm-store"

WORKDIR /usr/src/app
