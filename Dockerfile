# Étape 1 : Build
FROM golang:1.23.5-alpine3.21 AS builder

# Installer les dépendances nécessaires
RUN apk add --no-cache git

# Définir le répertoire de travail
WORKDIR /app

# Copier uniquement les fichiers de dépendances pour tirer parti du cache Docker
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste des fichiers sources
COPY . .

# Compiler l'application
RUN go build -o blockchain cmd/main.go

# Étape 2 : Image finale
FROM alpine:latest

# Définir le répertoire de travail
WORKDIR /app

# Copier le binaire depuis l'étape de build
COPY --from=builder /app/blockchain .

# Exposer le port 8000
EXPOSE 8000

# Commande par défaut
CMD ["./blockchain"]
