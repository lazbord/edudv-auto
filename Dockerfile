FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Set GOOS and GOARCH to ensure compatibility
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app with CGO disabled to create a static binary
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest

# Install tzdata for timezone support
RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /app/main .

COPY .env . 

# Set the timezone to Europe/Paris
ENV TZ=Europe/Paris

# Create a symlink to the correct timezone
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN chmod +x ./main

CMD ["./main"]
