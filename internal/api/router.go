package api

import (
	cust "semesta-ban/internal/api/customers"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/repository/repo_customers"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	mdl "github.com/go-chi/chi/middleware"
)

type ServerConfig struct {
	EncKey string
	JWTKey string
}

//todo add rate limiter
//todo add expired token from config
//implement credential email from config
func NewServer(db *sqlx.DB, cnf ServerConfig) *chi.Mux {
	var (
		r = chi.NewRouter()
		// ul = NewUnitLimiter()
		jwt    = localMdl.New([]byte(cnf.JWTKey))
		cuRepo = repo_customers.NewSqlRepository(db)
		cu     = cust.NewUsersHandler(db, cuRepo, jwt)
	)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	r.Use(mdl.RequestID)
	r.Use(mdl.RealIP)
	r.Use(mdl.Recoverer)

	r.Use(mdl.Heartbeat("/ping"))

	r.Post("/v1/login", cu.Login) //todo add auth for login & register endpoint
	r.Post("/v1/register", cu.Register)
	r.Get("/v1/verify", cu.VerifyEmail)

	r.Route("/v1/users", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Get("/me", cu.GetProfile)
	})

	return r
}
