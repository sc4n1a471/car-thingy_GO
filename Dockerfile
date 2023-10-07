# Latest golang image on apline linux
FROM golang:1.19-alpine

# Work directory
WORKDIR /docker-go

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "main.go"]

# Exposing server port
EXPOSE 3000