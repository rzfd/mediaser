# Frontend Dockerfile
# Multi-stage build for React application

# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Build arguments for React environment variables
ARG REACT_APP_GOOGLE_CLIENT_ID
ARG REACT_APP_API_URL=http://localhost:8080/api
ARG REACT_APP_BACKEND_URL=http://localhost:8080
ARG REACT_APP_CURRENCY_API=http://localhost:8094
ARG REACT_APP_LANGUAGE_API=http://localhost:8095

# Set environment variables for React build
ENV REACT_APP_GOOGLE_CLIENT_ID=$REACT_APP_GOOGLE_CLIENT_ID
ENV REACT_APP_API_URL=$REACT_APP_API_URL
ENV REACT_APP_BACKEND_URL=$REACT_APP_BACKEND_URL
ENV REACT_APP_CURRENCY_API=$REACT_APP_CURRENCY_API
ENV REACT_APP_LANGUAGE_API=$REACT_APP_LANGUAGE_API

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install --only=production

# Copy source code
COPY . .

# Build the React app with environment variables
RUN npm run build

# Production stage
FROM nginx:alpine

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy build files from React app
COPY --from=build /app/build /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 