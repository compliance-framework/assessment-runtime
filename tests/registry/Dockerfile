# Use the official nginx image as the base
FROM nginx:latest

# Remove the default nginx HTML files
RUN rm -rf /usr/share/nginx/html/*

# Copy the binary file to the nginx directory
COPY sample /usr/share/nginx/html/sample/1.0.0/sample

# Set permissions to make it readable
RUN chmod 755 /usr/share/nginx/html/sample/1.0.0/sample

# COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80
