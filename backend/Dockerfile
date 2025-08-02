# ====================================================
# FILE: backend/Dockerfile
# This is a multi-stage build to create a small, secure production image.
# ====================================================

# --- Build Stage ---
    FROM golang:1.24-alpine AS builder

    WORKDIR /app
    
    COPY . .
    
    RUN go mod tidy
    
    RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .
    
    # --- Final Stage ---
    FROM scratch
    
    # THIS IS THE FIX: Copy the schema file into the final image
    COPY --from=builder /app/schema.sql /schema.sql
    
    # Copy the built binary from the 'builder' stage
    COPY --from=builder /app/server /server
    
    EXPOSE 8080
    
    CMD ["/server"]