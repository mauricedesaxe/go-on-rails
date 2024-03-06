# Start with a base image that has both Go and Node.js
FROM golang:1.21.3-alpine

# Install system dependencies required for CGO and Node.js
RUN apk add --update nodejs npm gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy package.json and package-lock.json for Node.js dependencies
COPY package.json ./

# Install Node.js dependencies
RUN npm install

# Copy the rest of your source code
COPY . .

# Build your frontend assets
RUN npm run build:css

# Build your Go application
RUN go build -o myapp cmd/main.go

# Make port 3000 available to the world outside this container
EXPOSE 3000

# Run the executable
CMD ["./myapp"]