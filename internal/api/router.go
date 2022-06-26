package api

import (
	"semesta-ban/internal/api/auth"
	cust "semesta-ban/internal/api/customers"
	"semesta-ban/internal/api/master_data"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/internal/api/products"
	"semesta-ban/repository/repo_customers"
	"semesta-ban/repository/repo_master_data"
	"semesta-ban/repository/repo_products"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	mdl "github.com/go-chi/chi/middleware"
)

type ServerConfig struct {
	EncKey            string
	JWTKey            string
	AnonymousKey      string
	BaseAssetsUrl     string
	UploadPath        string
	ProfilePicPath    string
	ProfilePicMaxSize int
}

//todo add rate limiter
//todo add expired token from config
//todo add base url from config for profile picture , product picture
//implement credential email from config
func NewServer(db *sqlx.DB, cnf ServerConfig) *chi.Mux {
	var (
		r = chi.NewRouter()
		// ul = NewUnitLimiter()
		jwt               = localMdl.New([]byte(cnf.JWTKey))
		anon              = localMdl.New([]byte(cnf.AnonymousKey))
		cuRepo            = repo_customers.NewSqlRepository(db)
		prRepo            = repo_products.NewSqlRepository(db)
		mdRepo            = repo_master_data.NewSqlRepository(db)
		custHandler       = cust.NewUsersHandler(db, cuRepo, jwt, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.ProfilePicPath, cnf.ProfilePicMaxSize)
		authHandler       = auth.NewAuthHandler(jwt, anon)
		prodHandler       = products.NewProductsHandler(db, prRepo, cnf.BaseAssetsUrl)
		masterDataHandler = master_data.NewMasterDataHandler(db, mdRepo, cnf.BaseAssetsUrl)
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

	r.Get("/v1/verify", custHandler.VerifyEmail)
	r.Get("/v1/auth/anonymous", authHandler.GetAnonymousToken)

	r.Route("/v1", func(r chi.Router) { //anonymous scope
		r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
		r.Post("/login", custHandler.Login)
		r.Post("/register", custHandler.Register)
		r.Route("/master-data", func(r chi.Router) {
			r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
			r.Get("/tire-brand", masterDataHandler.GetListMerkBan)
			r.Get("/gender", masterDataHandler.GetListGender)
			r.Get("/outlet", masterDataHandler.GetListOutlet)
			// r.Get("/motor-brand", prodHandler.GetListProducts)
			// r.Get("/outlets", prodHandler.GetListProducts)
		})

	})

	r.Route("/v1/auth", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Get("/refresh", authHandler.RefreshToken)
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Get("/me", custHandler.GetProfile)
		r.Post("/change-password", custHandler.ChangePassword)
		r.Post("/resend-email", custHandler.ResendEmailVerification)
		r.Post("/request-pin-email", custHandler.RequestPinEmail)
		r.Post("/change-email", custHandler.ChangeEmail)
		r.Post("/update-name", custHandler.UpdateName)
		r.Post("/update-gender", custHandler.UpdateGender)
		r.Post("/update-phone", custHandler.UpdatePhoneNumber)
		r.Post("/update-birthdate", custHandler.UpdateBirthDate)
		r.Post("/upload-profile-img", custHandler.UploadProfileImg)
	})

	r.Route("/v1/products", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
		r.Get("/", prodHandler.GetListProducts)
		r.Get("/detail", prodHandler.GetProductDetail)
	})

	r.Route("/v1/products/cart", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/add", prodHandler.CartAdd)
		r.Post("/remove", prodHandler.CartRemove)
		r.Post("/update", prodHandler.CartUpdate)
		r.Post("/select-deselect-all", prodHandler.CartSelectDeselectAll)
		r.Get("/me", prodHandler.CartMe)
	})

	r.Route("/v1/products/wishlist", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/add", prodHandler.WishlistAdd)
		r.Post("/remove", prodHandler.WishlistRemove)
		r.Get("/me", prodHandler.WishlistMe)
	})

	return r
}
