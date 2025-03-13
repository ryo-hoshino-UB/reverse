FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN npm install -g pnpm@10.4.1
WORKDIR /app

FROM base as builder
COPY package.json pnpm-lock.yaml* ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm run -r build:prod

# 本番環境用ビルド
FROM builder as othello-web
WORKDIR /app
ENV NODE_ENV=production
EXPOSE 5001
CMD [ "pnpm", "start"]