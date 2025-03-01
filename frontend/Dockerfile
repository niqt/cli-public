# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copy package.json e package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

COPY . .

# Build
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy the file in the nginx path
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
