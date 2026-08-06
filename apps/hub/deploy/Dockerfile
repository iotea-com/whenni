# --- Build Stage ---
FROM node:24-alpine AS builder

# Set working directory
WORKDIR /build

# Install dependencies
COPY package.json pnpm-lock.yaml* pnpm-workspace.yaml nx.json migrations.json tsconfig.base.json .eslintrc.json ./
COPY ./prisma ./prisma
RUN corepack enable && corepack prepare pnpm@11.16.0 --activate
RUN pnpm install --frozen-lockfile
RUN pnpm exec prisma generate --generator client

# Copy source code
COPY ./libs/frontend ./libs/frontend
COPY ./libs/id ./libs/id
COPY ./libs/gruent-js ./libs/gruent-js
COPY ./libs/http ./libs/http
COPY ./libs/engine/channels ./libs/engine/channels
COPY ./libs/engine/nodes/v1 ./libs/engine/nodes/v1
COPY ./libs/engine/dependencies ./libs/engine/dependencies
COPY ./apps/hub ./apps/hub

# Build the Next.js app
RUN pnpm nx build hub --configuration=production

# --- Production Stage ---
FROM node:24-alpine AS runner

WORKDIR /hub

# Set NODE_ENV to production
ENV NODE_ENV=production

# Configure Next.js to bind to 0.0.0.0
ENV HOST=0.0.0.0
ENV PORT=3000

# Copy and install dependencies
RUN corepack enable && corepack prepare pnpm@11.16.0 --activate
COPY --from=builder /build/pnpm-lock.yaml ./
COPY --from=builder /build/pnpm-workspace.yaml ./
COPY --from=builder /build/apps/hub/dist .
COPY --from=builder /build/prisma ./prisma
RUN pnpm install --frozen-lockfile
COPY --from=builder /build/node_modules/@prisma ./node_modules/@prisma
COPY --from=builder /build/node_modules/prisma ./node_modules/prisma
RUN pnpm dlx prisma generate --generator client

# Expose the port
EXPOSE 3000

# Start the app
CMD ["pnpm", "start"]