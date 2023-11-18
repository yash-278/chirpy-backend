module chirpy-backend

go 1.20

require github.com/go-chi/chi/v5 v5.0.10

require github.com/yash-278/chirpy-backend/database v1.0.0

require golang.org/x/crypto v0.15.0 // indirect

replace github.com/yash-278/chirpy-backend/database => ./internal/database
