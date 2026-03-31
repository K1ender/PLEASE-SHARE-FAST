# PSF - Please Share Fast

A lightweight file sharing service built with Go. Upload files and share them with a simple link.

## Features

- File upload and download via HTTP
- In-memory file storage with disk persistence
- Automatic cleanup of files older than 24 hours
- Request logging and request ID tracking
- Configurable via environment variables
- Docker support

## Requirements

- Go 1.26.1 or later
- Docker (optional)

## Installation

Clone the repository:

```bash
git clone https://github.com/k1ender/psf.git
cd psf
```

## Configuration

Create a `.env` file based on `.env.example`:

```env
ENV=production
HTTP_PORT=8080
HTTP_ADDR=0.0.0.0
```

| Variable   | Default    | Description              |
|------------|------------|--------------------------|
| `ENV`      | production | Environment mode         |
| `HTTP_PORT`| 8080       | Server port              |
| `HTTP_ADDR`| 0.0.0.0    | Server address           |

## Usage

### Running locally

```bash
go run cmd/api/main.go
```

Or using Task:

```bash
task
```

### Running with Docker

```bash
docker build -t psf .
docker run -p 8080:8080 psf
```

## API Endpoints

| Method | Path         | Description              |
|--------|--------------|--------------------------|
| GET    | `/`          | Upload page              |
| POST   | `/upload`    | Upload a file            |
| GET    | `/file/{id}` | Download file by ID      |

### Upload a file

```bash
curl -X POST -F "file=@myfile.txt" http://localhost:8080/upload
```

Response: Returns the file ID (e.g., `abc123`)

### Download a file

```bash
curl -O http://localhost:8080/file/abc123
```

## Architecture

```
psf/
├── cmd/api/           # Application entry point
├── internal/
│   ├── cleaner/       # File cleanup logic
│   ├── config/        # Configuration management
│   ├── logger/        # Logging utilities
│   ├── middleware/    # HTTP middleware
│   ├── model/         # Data models
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── transport/http/# HTTP handlers
├── pkg/api/           # Main application runner
├── templates/         # HTML templates
├── data/              # File storage directory
└── Dockerfile
```
