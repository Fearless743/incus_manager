module incus-manager

go 1.21

require (
	github.com/gorilla/websocket v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.17.0
)

require github.com/golang-jwt/jwt/v5 v5.3.1 // indirect

replace github.com/lib/pq => github.com/lib/pq v1.10.9
