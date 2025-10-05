# ====== 1️⃣ Build Stage ======
FROM debian:bookworm-slim AS build

# Install dependencies for Flutter
RUN apt-get update && apt-get install -y \
    curl git unzip xz-utils zip libglu1-mesa && \
    rm -rf /var/lib/apt/lists/*

# Set Flutter version (you can change it if needed)
ENV FLUTTER_VERSION=3.35.5

# Install Flutter SDK
RUN git clone https://github.com/flutter/flutter.git /flutter -b $FLUTTER_VERSION

# Add Flutter to PATH
ENV PATH="/flutter/bin:/flutter/bin/cache/dart-sdk/bin:${PATH}"

# Pre-download Flutter dependencies
RUN flutter doctor

# Copy project files
WORKDIR /app
COPY . .

# Enable web support & get dependencies
RUN flutter config --enable-web && flutter pub get

# Build the web app (release mode)
RUN flutter build web --release

# ====== 2️⃣ Runtime Stage ======
FROM nginx:stable-alpine

# Remove default nginx web content
RUN rm -rf /usr/share/nginx/html/*

# Copy Flutter build output to nginx html folder
COPY --from=build /app/build/web /usr/share/nginx/html

# Copy custom nginx config (optional)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
