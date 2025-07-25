ARG IMAGE_UV
FROM ${IMAGE_UV}

ENV UV_COMPILE_BYTECODE=1
ENV UV_LINK_MODE=copy
ENV UV_PROJECT_ENVIRONMENT=/usr/local

WORKDIR /db

# Install the project's dependencies using the lockfile and settings
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=db/uv.lock,target=uv.lock \
    --mount=type=bind,source=db/pyproject.toml,target=pyproject.toml \
    uv sync --locked --no-dev

# Reset the entrypoint, don't invoke `uv`
ENTRYPOINT []
