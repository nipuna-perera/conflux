# Multi-stage Docker build for SvelteKit frontend
# Stage 1: Build the SvelteKit application
FROM node:24-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .

# Set environment variables for build time
ARG PUBLIC_API_URL=http://localhost:8080
ENV PUBLIC_API_URL=$PUBLIC_API_URL

RUN npm run build

# Stage 2: Serve the built application
FROM node:24-alpine
WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm ci --only=production

EXPOSE 3000
CMD ["node", "build"]
