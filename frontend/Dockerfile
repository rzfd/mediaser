# Frontend Dockerfile
# Multi-stage build for efficient static file serving

# Stage 1: Development stage with Node.js (optional for build tools)
FROM node:18-alpine AS development

WORKDIR /app

# Copy package files if they exist (for future expansion)
COPY package*.json ./

# Install dependencies if package.json exists
RUN if [ -f package.json ]; then npm install; fi

# Copy frontend files
COPY . .

# Stage 2: Production stage with Nginx
FROM nginx:alpine AS production

# Install additional tools for debugging
RUN apk add --no-cache curl wget

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy frontend files to nginx html directory
COPY --from=development /app/*.html /usr/share/nginx/html/
COPY --from=development /app/*.js /usr/share/nginx/html/
COPY --from=development /app/*.css /usr/share/nginx/html/
COPY --from=development /app/*.md /usr/share/nginx/html/

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

# Create nginx user and set permissions
RUN addgroup -g 1001 -S nginx-user && \
    adduser -S nginx-user -u 1001 -G nginx-user && \
    chown -R nginx-user:nginx-user /usr/share/nginx/html && \
    chown -R nginx-user:nginx-user /var/cache/nginx && \
    chown -R nginx-user:nginx-user /var/log/nginx && \
    chown -R nginx-user:nginx-user /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R nginx-user:nginx-user /var/run/nginx.pid

# Switch to non-root user
USER nginx-user

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 