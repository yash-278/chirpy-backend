module chirpy-backend

go 1.20

require github.com/go-chi/chi/v5 v5.0.10

require github.com/yash-278/chirpy-backend/database v1.0.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/yash-278/chirpy-backend/auth v1.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.1.0
	golang.org/x/crypto v0.17.0 // indirect
)

replace github.com/yash-278/chirpy-backend/database => ./internal/database

replace github.com/yash-278/chirpy-backend/auth => ./internal/auth
