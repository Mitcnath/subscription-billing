# Use the official nginx Alpine image to serve static files.
# No build step needed — the frontend is plain HTML.
FROM nginx:alpine

# Copy the HTML file into nginx's default serve directory.
# nginx will serve this automatically at http://localhost:80.
COPY index.html /usr/share/nginx/html/index.html

# Document that nginx listens on port 80.
EXPOSE 80