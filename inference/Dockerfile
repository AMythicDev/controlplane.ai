FROM ghcr.io/astral-sh/uv:debian

WORKDIR /app

ENV UV_COMPILE_BYTECODE=1

COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-cache --no-install-project
COPY model/ ./model/
COPY src ./src/
COPY README.md ./
RUN uv sync --frozen --no-cache
CMD ["uv", "run", "src/inference/main.py"]
