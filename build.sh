#!/bin/bash
set -e

echo "=== Building Incus Manager ==="

# Build frontend
echo "[1/3] Building frontend..."
cd frontend
npm ci --silent 2>/dev/null || npm install --silent
npm run build
cd ..
echo "  ✅ Frontend built"

# Copy frontend dist to backend
echo "[2/3] Preparing backend..."
rm -rf backend/dist
cp -r frontend/dist backend/dist
echo "  ✅ Frontend assets copied"

# Build backend binary
echo "[3/3] Building backend..."
cd backend
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ../incus-manager cmd/server.go
cd ..
chmod +x incus-manager
echo "  ✅ Backend built"

echo ""
echo "=== Build Complete ==="
echo ""
echo "Single binary deployment:"
echo "  ./incus-manager"
echo ""
echo "Or use Docker Compose (single port 80):"
echo "  docker-compose up --build -d"
echo ""
echo "Or local development:"
echo "  cd frontend && npm run dev     # frontend on :5173"
echo "  cd backend && go run cmd/server.go  # backend on :8080"
