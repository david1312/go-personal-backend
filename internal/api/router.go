package api

import (
	"semesta-ban/internal/api/auth"
	cust "semesta-ban/internal/api/customers"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/internal/api/products"
	"semesta-ban/repository/repo_customers"
	"semesta-ban/repository/repo_products"

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
//todo add base url from config for profile picture , product picture
//implement credential email from config
func NewServer(db *sqlx.DB, cnf ServerConfig) *chi.Mux {
	var (
		r = chi.NewRouter()
		// ul = NewUnitLimiter()
		jwt         = localMdl.New([]byte(cnf.JWTKey))
		cuRepo      = repo_customers.NewSqlRepository(db)
		prRepo      = repo_products.NewSqlRepository(db)
		custHandler = cust.NewUsersHandler(db, cuRepo, jwt)
		authHandler = auth.NewAuthHandler(jwt)
		prodHandler = products.NewProductsHandler(db, prRepo)
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

	r.Post("/v1/login", custHandler.Login) //todo add auth for login & register endpoint
	r.Post("/v1/register", custHandler.Register)
	r.Get("/v1/verify", custHandler.VerifyEmail)

	r.Route("/v1/auth", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Get("/refresh", authHandler.RefreshToken)
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Get("/me", custHandler.GetProfile)
		r.Post("/change-password", custHandler.ChangePassword)
		r.Post("/resend-email", custHandler.ResendEmailVerification)
		r.Post("/request-pin-email", custHandler.RequestPinEmail)
		r.Post("/change-email", custHandler.ChangeEmail)
	})

	r.Route("/v1/products", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware)
		r.Get("/", prodHandler.GetListProducts)
	})

	return r
}
