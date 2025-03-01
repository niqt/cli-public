# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copia package.json e package-lock.json
COPY package*.json ./

# Installa le dipendenze
RUN npm install

# Copia il codice sorgente
COPY . .

# Compila l'applicazione
RUN npm run build

# Production stage
FROM nginx:alpine

# Copia i file compilati nella directory di Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Copia la configurazione di Nginx personalizzata (opzionale)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Espone la porta 80
EXPOSE 80

# Avvia Nginx
CMD ["nginx", "-g", "daemon off;"]
