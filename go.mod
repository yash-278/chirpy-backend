module chirpy-backend

go 1.20

require github.com/go-chi/chi/v5 v5.0.10

require github.com/yash-278/chirpy-backend/database v1.0.0

replace github.com/yash-278/chirpy-backend/database => ./internal/database
