# Use specific versions for both Go and Node.js to ensure reproducible builds.
# Alpine is used for its small size and package availability.
FROM golang:1.22-alpine AS go-builder

# Install build dependencies for Go.
# No additional packages are needed at this stage, so the RUN command for apk add is removed.

WORKDIR /app

# Copy the go.mod and go.sum files separately to cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source code.
COPY . ./

# https://gin-gonic.com/docs/jsoniter/
RUN go build -o go-htmx -tags=jsoniter

# Use a separate stage for Node.js to keep the build stages focused and reusable.
FROM node:16-alpine AS node-builder

WORKDIR /nodeapp

# Copy package.json and package-lock.json for npm install.
# This is done before copying the rest of the code to leverage Docker cache.
COPY package-lock.json package.json ./
RUN npm install --save

# Copy the necessary files for building the CSS with Tailwind.
# Assume that only specific directories or files are needed for this step.
# Adjust the paths if your static files are located elsewhere.
COPY . .

# Use npx to run Tailwind CLI without global install.
# The input and output paths are corrected to match the current directory structure.
RUN npx tailwindcss -i './web/static/css/tailwind.css' -o './web/static/css/styles.css'

# Final stage: Combine the Go build and Node.js build in the final image.
FROM alpine:latest

# Copy the Go build output from the Go build stage.
# Adjust the copy command according to your Go build output.
COPY --from=go-builder /app/go-htmx /app/go-htmx

# Copy the static files with compiled CSS from the Node.js build stage.
COPY --from=node-builder /nodeapp/web /app/web

# Set the working directory to /app to align with the copied directory structure.
WORKDIR /app
