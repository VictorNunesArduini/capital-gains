FROM golang:1.23.2-alpine AS build

# Set the working directory in the container
WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

# Build the executable
RUN go build -o main ./cmd

# BUILD
FROM alpine:latest

RUN apk update

WORKDIR /app/capital-gains

COPY --from=build /app/main .
COPY --from=build /app/input.txt .

RUN echo "start running capital-gains app"

# Command to run the app
ENTRYPOINT ["./main"]