FROM node:22-alpine AS base

RUN npm install -g npm@11.14.1

FROM base AS deps

WORKDIR /app

ARG NEXT_PUBLIC_API_BASE_URL=""
ARG INTERNAL_API_BASE_URL="http://portlyn:8080"
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV INTERNAL_API_BASE_URL=$INTERNAL_API_BASE_URL

COPY package.json package-lock.json ./
RUN npm install

FROM base AS builder

WORKDIR /app

ARG NEXT_PUBLIC_API_BASE_URL=""
ARG INTERNAL_API_BASE_URL="http://portlyn:8080"
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV INTERNAL_API_BASE_URL=$INTERNAL_API_BASE_URL

COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npm run build

FROM node:22-alpine

WORKDIR /app

ARG NEXT_PUBLIC_API_BASE_URL=""
ARG INTERNAL_API_BASE_URL="http://portlyn:8080"
ENV NODE_ENV=production
ENV PORT=3000
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV INTERNAL_API_BASE_URL=$INTERNAL_API_BASE_URL

RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/bin/npm /usr/local/bin/npx

COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/package-lock.json ./package-lock.json
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.mjs ./next.config.mjs

EXPOSE 3000

CMD ["node", "node_modules/next/dist/bin/next", "start"]
