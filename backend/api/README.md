# API Specifications

This directory contains OpenAPI/Swagger specification files for the API.

## Structure

- `openapi.yaml` - OpenAPI 3.0 specification
- `swagger.json` - Generated Swagger file

## Generate Swagger Documentation

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
cd backend
swag init -g cmd/api/main.go -o swagger

```

