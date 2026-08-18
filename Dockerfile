FROM oven/bun:1.3.14-alpine@sha256:5acc90a93e91ff07bf72aa90a7c9f0fa189765aec90b47bdbf2152d2196383c0

WORKDIR /app
ENV NODE_ENV=production

COPY package.json tsconfig.json bun.lock bunfig.toml ./
RUN bun install --frozen-lockfile --production --ignore-scripts
COPY src/ ./src/
USER bun

CMD ["bun", "src/main.ts"]
