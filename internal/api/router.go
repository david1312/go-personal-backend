package api

import (
	"libra-internal/internal/api/auth"
	cust "libra-internal/internal/api/customers"
	"libra-internal/internal/api/master_data"
	"libra-internal/internal/api/merchant"
	localMdl "libra-internal/internal/api/middleware"
	"libra-internal/internal/api/products"
	"libra-internal/internal/api/ratings"
	"libra-internal/internal/api/reports"
	"libra-internal/internal/api/transactions"
	"libra-internal/repository/repo_customers"
	"libra-internal/repository/repo_master_data"
	"libra-internal/repository/repo_merchant"
	"libra-internal/repository/repo_products"
	"libra-internal/repository/repo_ratings"
	"libra-internal/repository/repo_reports"
	"libra-internal/repository/repo_transactions"
	"net/http"

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
	MaxFileSize       int
	ProfilePicPath    string
	ProfilePicMaxSize int
	MidtransConfig    transactions.MidtransConfig
}

// todo add rate limiter
// todo add expired token from config
// todo add base url from config for profile picture , product picture
// implement credential email from config
func NewServer(db *sqlx.DB, client *http.Client, cnf ServerConfig) *chi.Mux {
	var (
		r = chi.NewRouter()
		// ul = NewUnitLimiter()
		jwt               = localMdl.New([]byte(cnf.JWTKey))
		anon              = localMdl.New([]byte(cnf.AnonymousKey))
		cuRepo            = repo_customers.NewSqlRepository(db)
		prRepo            = repo_products.NewSqlRepository(db)
		mdRepo            = repo_master_data.NewSqlRepository(db)
		trRepo            = repo_transactions.NewSqlRepository(db)
		rateRepo          = repo_ratings.NewSqlRepository(db)
		reportsRepo       = repo_reports.NewSqlRepository(db)
		custHandler       = cust.NewUsersHandler(db, cuRepo, jwt, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.ProfilePicPath, cnf.ProfilePicMaxSize)
		authHandler       = auth.NewAuthHandler(jwt, anon)
		prodHandler       = products.NewProductsHandler(db, prRepo, mdRepo, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.MaxFileSize)
		transHandler      = transactions.NewTransactionsHandler(db, prRepo, mdRepo, trRepo, cnf.BaseAssetsUrl, client, cnf.MidtransConfig)
		masterDataHandler = master_data.NewMasterDataHandler(db, mdRepo, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.MaxFileSize)
		rateHandler       = ratings.NewRatingsHandler(db, rateRepo, prRepo, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.MaxFileSize)
		reportHandler     = reports.NewReportsHandler(reportsRepo)

		//merchant
		merchRepo       = repo_merchant.NewSqlRepository(db)
		merchantHandler = merchant.NewMerchantHandler(merchRepo, prRepo, jwt, cnf.BaseAssetsUrl, cnf.UploadPath, cnf.ProfilePicMaxSize)
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
	r.Get("/version", authHandler.GetVersion)

	r.Route("/v1", func(r chi.Router) { //anonymous scope
		r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
		r.Post("/login", custHandler.Login)
		r.Post("/signin-google", custHandler.SignInGoogle)
		r.Post("/register", custHandler.Register)
		r.Post("/check-location", masterDataHandler.CheckLocation)
		r.Route("/master-data", func(r chi.Router) {
			r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
			r.Get("/tire-brand", masterDataHandler.GetListMerkBan)
			r.Get("/gender", masterDataHandler.GetListGender)
			r.Get("/outlet", masterDataHandler.GetListOutlet)
			r.Get("/sort-by", masterDataHandler.GetListSortBy)
			r.Get("/tire-size", masterDataHandler.GetListSizeBan)
			r.Get("/tire-size-raw", masterDataHandler.EPGetListTireSizeRaw)
			r.Get("/motor-brand", masterDataHandler.GetListMerkMotor)
			r.Get("/motor-list-by-brand", masterDataHandler.GetListMotorByBrand)
			r.Get("/payment-method", masterDataHandler.GetListPaymentMethod)
			r.Get("/toprank-motor", masterDataHandler.GetTopRankMotor)
			r.Get("/asset-img", masterDataHandler.GetImgAsset)
			r.Get("/tire-type", masterDataHandler.GetTireType)
			r.Get("/magic", masterDataHandler.MagicHandler)
			r.Get("/magic-transaction-expired", masterDataHandler.TransactionExpiredHandler)
			// r.Get("/outlets", prodHandler.GetListProducts)
		})

	})

	r.Route("/v1/auth", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Get("/refresh", authHandler.RefreshToken)
	})

	r.Route("/v1/callback", func(r chi.Router) {
		r.Post("/midtrans-payment", transHandler.CallbackPayment)
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

	r.Route("/v1/transactions", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/submit", transHandler.SubmitTransactions)
		r.Get("/inquiry/schedule", transHandler.InquirySchedule)
		r.Post("/history", transHandler.GetHistoryTransactions)
		r.Post("/payment-instruction", transHandler.GetPaymentInstruction)
		r.Post("/detail", transHandler.GetTransactionDetail)
		r.Get("/count", transHandler.GetCountTransaction)
		r.Get("/test", transHandler.TestJubelio)
	})

	r.Route("/v1/products", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAnonymous))
		r.Post("/", prodHandler.GetListProducts)
		r.Get("/detail", prodHandler.GetProductDetail)
	})

	r.Route("/v1/products/cart", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/add", prodHandler.CartAdd)
		r.Post("/remove", prodHandler.CartRemove)
		r.Post("/update", prodHandler.CartUpdate)
		r.Post("/select-deselect-all", prodHandler.CartSelectDeselectAll)
		r.Get("/me", prodHandler.CartMe)
		r.Get("/summary", prodHandler.GetCartSummary)
	})

	r.Route("/v1/products/wishlist", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/add", prodHandler.WishlistAdd)
		r.Post("/remove", prodHandler.WishlistRemove)
		r.Get("/me", prodHandler.WishlistMe)
	})

	r.Route("/v1/ratings", func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(localMdl.GuardAccess))
		r.Post("/product/submit", rateHandler.SubmitRatingProduct)
		r.Post("/outlet/submit", rateHandler.SubmitRatingOutlet)
		r.Post("/outlet/history", rateHandler.GetListRatingOutler)
	})

	r.Route("/v1/merchant/auth", func(r chi.Router) { //anonymous scope
		r.Use(jwt.AuthMiddlewareMerchant(localMdl.GuardAnonymous))
		r.Post("/login", merchantHandler.LoginMerchant)

	})

	r.Route("/v1/merchant", func(r chi.Router) { //anonymous scope
		r.Use(jwt.AuthMiddlewareMerchant(localMdl.GuardAccess))
		r.Get("/me", merchantHandler.GetProfileMerchant)

		r.Route("/reports", func(r chi.Router) {
			r.Post("/sales/syncup", reportHandler.SyncSales)
			r.Post("/sales/calculate-profit", reportHandler.SalesCalculateProfit)
			r.Post("/sales/get/all", reportHandler.EPGetSalesReport)
		})

		r.Route("/products", func(r chi.Router) {
			r.Post("/delete", prodHandler.DeleteProduct)
			r.Post("/add", prodHandler.AddProduct)
			r.Post("/update", prodHandler.EPProductUpdate)
			r.Post("/images-add", prodHandler.EPProductAddImage)
			r.Post("/images-delete", prodHandler.EpProductDeleteImage)
			r.Post("/images-update", prodHandler.EpProductUpdateImage)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Post("/update-status", transHandler.EPUpdateTransactionStatus)
			r.Post("/history", transHandler.EPMerchantGetHistoryTransactions)
			r.Post("/detail", transHandler.EPMerchantGetTransactionDetail)
		})

		r.Route("/master-data", func(r chi.Router) {
			r.Post("/brand-motor/add", masterDataHandler.EPAddBrandMotor)
			r.Post("/brand-motor/delete", masterDataHandler.EPRemoveBrandMotor)
			r.Post("/brand-motor/update", masterDataHandler.EPUpdateBrandMotor)
			r.Post("/brand-motor/update-image", masterDataHandler.EPUpdateBrandMotorIcon)

			r.Post("/tire-brand/add", masterDataHandler.EPAddTireBrand)
			r.Post("/tire-brand/delete", masterDataHandler.EPRemoveTireBrand)
			r.Post("/tire-brand/update", masterDataHandler.EPUpdateTireBrand)
			r.Post("/tire-brand/update-image", masterDataHandler.EPUpdateTireBrandIcon)

			r.Post("/motor", masterDataHandler.EPListMotor)
			r.Post("/motor/add", masterDataHandler.EPMotorAdd)
			r.Post("/motor/update", masterDataHandler.EPMotorUpdate)
			r.Post("/motor/update-image", masterDataHandler.EPMotorUpdateImage)
			r.Post("/motor/delete", masterDataHandler.EPMotorRemove)

			r.Get("/category-motor", masterDataHandler.EPCategoryMotor)

			r.Post("/tire-size/add", masterDataHandler.EPTireSizeAdd)
			r.Post("/tire-size/delete", masterDataHandler.EPTireSizeDelete)

			r.Post("/tire-ring/add", masterDataHandler.EPTireRingAdd)

		})

	})

	return r
}
