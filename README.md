# Go Utils Library

## Overview
This utility library provides common functionality for logging, Firebase integration, and database connections in Go projects. It simplifies common tasks and provides a standardized way to handle these functionalities across your Go applications.

## Requirements
- Go 1.24.1 or higher
- Firebase Admin SDK (for Firebase utilities)
- Database drivers (depending on your database implementation)

## Installation

```$ go get github.com/siva-chegondi/go-utils```

## Features

### Zerolog Logging Utilities
- `DefaultLogger` with zerolog logging setup and configuration
- Log level management with `LOG_LEVEL` environment variable.
- Logger Middleware for gin framework

### Firebase Integration
- Firebase initialization and configuration.
- Authentication `VerifyToken` token verification helper.

### Database Utilities
- Connection pool management
- Database connection helpers
- Transaction utilities
- Query builders
