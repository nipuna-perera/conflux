# Multi-stage Docker build for SvelteKit frontend
# Stage 1: Build the SvelteKit application
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Stage 2: Serve the built application
FROM node:18-alpine
WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm ci --only=production

EXPOSE 3000
CMD ["node", "build"]
