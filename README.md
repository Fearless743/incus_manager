# Incus Manager

A web-based management panel for Incus container management with multi-user support, host management, and instance sharing capabilities.

## Features

- Multi-user authentication with JWT
- Add unlimited hosts per user
- Instance creation with custom configurations
- Automatic IP allocation based on host
- Port mapping configuration
- Resource limits (CPU, Memory, Disk, Network)
- Instance sharing between users with expiry dates
- Real-time monitoring via WebSocket
- Project-based isolation using Incus projects

## Tech Stack

### Backend
- Go 1.21+
- GORM (ORM)
- PostgreSQL
- JWT Authentication
- WebSocket for real-time updates
- Incus REST API

### Frontend
- React 18
- Vite
- React Router
- Axios
- WebSocket

## Project Structure

```
incus_manager/
├── backend/
│   ├── cmd/
│   │   └── server.go          # Main entry point
│   ├── internal/
│   │   ├── config/            # Configuration
│   │   ├── handler/           # HTTP handlers
│   │   ├── middleware/        # Auth, CORS, Logging
│   │   ├── model/             # Data models
│   │   ├── service/           # Business logic
│   │   └── websocket/         # WebSocket hub
│   └── pkg/
│       └── database/          # Database setup
├── frontend/
│   └── src/
│       ├── components/        # React components
│       ├── context/           # Auth context
│       ├── pages/             # Page components
│       └── services/          # API calls
└── docker-compose.yml
```

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Incus installed and configured

## Setup

### Backend

1. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

2. Update the environment variables in `.env`

3. Install dependencies and run:
```bash
cd backend
go mod tidy
go run cmd/server.go
```

### Frontend

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Run development server:
```bash
npm run dev
```

### Docker

```bash
docker-compose up -d
```

## API Endpoints

See `API_DOCUMENTATION.txt` for complete API documentation.

## Incus Integration

The system uses Incus projects to isolate different hosts' instances. Each host gets its own project with automatic IP allocation.

## License

MIT
