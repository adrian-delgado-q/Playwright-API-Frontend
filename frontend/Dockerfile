# Use nginx alpine for lightweight frontend serving
FROM nginx:alpine

# Install curl for health checks
RUN apk add --no-cache curl

# Copy the HTML file to nginx html directory
COPY books_demo.html /usr/share/nginx/html/index.html

# Copy our server configuration to replace default
COPY default.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://127.0.0.1/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
