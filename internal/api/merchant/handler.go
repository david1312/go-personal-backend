package merchant

import (
	"libra-internal/internal/api/customers"
	localMdl "libra-internal/internal/api/middleware"
	"libra-internal/internal/api/response"
	"libra-internal/pkg/constants"
	"libra-internal/pkg/crashy"
	"libra-internal/repository/repo_merchant"
	"libra-internal/repository/repo_products"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type MerchantHandler struct {
	merchRepo         repo_merchant.MerchantRepository
	prodRepo          repo_products.ProductsRepository
	jwt               *localMdl.JWT
	baseAssetUrl      string
	uploadPath        string
	profilePicMaxSize int
}

func NewMerchantHandler(mr repo_merchant.MerchantRepository, pr repo_products.ProductsRepository, jwt *localMdl.JWT, baseAssetUrl, uploadPath string, profilePicMaxSize int) *MerchantHandler {
	return &MerchantHandler{merchRepo: mr, prodRepo: pr, jwt: jwt, baseAssetUrl: baseAssetUrl, uploadPath: uploadPath, profilePicMaxSize: profilePicMaxSize}
}

func (mrc *MerchantHandler) LoginMerchant(w http.ResponseWriter, r *http.Request) {
	var (
		p   LoginRequest
		ctx = r.Context()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	merchant, errCode, err := mrc.merchRepo.Login(ctx, p.Username, p.Password)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//generate token
	expiredTime := time.Now().Add(24 * 7 * time.Hour)
	_, tokenLogin, _ := mrc.jwt.JWTAuth.Encode(&localMdl.MerchantToken{
		OutletId: merchant.OutletId,
		Username: merchant.Username,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 30)
	_, tokenRefresh, _ := mrc.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      merchant.Username,
		CustName: merchant.Username,
		Expired:  expiredTimeRefresh,
	})

	response.Yay(w, r, customers.LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
	}, http.StatusOK)
}

func (mrc *MerchantHandler) GetProfileMerchant(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.MerchantToken)
	)

	merchant, errCode, err := mrc.merchRepo.GetMerchantProfile(ctx, authData.Username)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, MerchantDataResponse{
		Username:      merchant.Username,
		OutletId:      merchant.OutletId,
		OutletName:    merchant.OutletName,
		OutletAvatar:  mrc.baseAssetUrl + constants.MerchantDir + merchant.OutletAvatar,
		OutletEmail:   merchant.OutletEmail,
		CsNumber:      merchant.CsNumber,
		OutletAddress: merchant.OutletAddress,
		OutletCity:    merchant.OutletCity,
		OutletGmapUrl: merchant.OutletGmapUrl,
	}, http.StatusOK)
}
