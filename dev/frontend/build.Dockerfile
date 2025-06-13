ARG IMAGE_NODE
ARG IMAGE_NGINX
FROM ${IMAGE_NODE} AS build

ARG REPO_NAME
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
RUN pnpm config set store-dir "$PNPM_HOME/.pnpm-store"

WORKDIR /usr/src/app/frontend
COPY frontend/package*.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm run build

FROM ${IMAGE_NGINX}
COPY --from=build /usr/src/app/frontend/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
