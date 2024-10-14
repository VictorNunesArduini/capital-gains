FROM golang:latest-alpine AS build

# Set the working directory in the container
WORKDIR /app

COPY cmd/main.go .

# Build the executable
RUN go build -o main .

# BUILD
FROM alpine:latest

RUN apk update

COPY --from=build /app/capital-gains .

WORKDIR /app/capital-gains

RUN echo "start running capital-gains app"

# Command to run the app
ENTRYPOINT ["./main"]

CMD ["input.txt"]