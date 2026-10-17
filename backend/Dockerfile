# Build stage: compile a static binary
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/godo ./cmd/godo

# Run stage: just the binary, as a non-root user
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/godo /godo
EXPOSE 8080
ENTRYPOINT ["/godo"]
