package products

import (
	"errors"
	"net/http"
	"semesta-ban/internal/api/response"
	cn "semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/repository/repo_products"
	"strconv"

	localMdl "semesta-ban/internal/api/middleware"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type ProductsHandler struct {
	db           *sqlx.DB
	prodRepo     repo_products.ProductsRepository
	baseAssetUrl string
}

func NewProductsHandler(db *sqlx.DB, pr repo_products.ProductsRepository, baseAssetUrl string) *ProductsHandler {
	return &ProductsHandler{db: db, prodRepo: pr, baseAssetUrl: baseAssetUrl}
}

func (prd *ProductsHandler) GetListProducts(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		fp  = NewProductsParams(r)
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if (len(fp.OrderBy) > 0 && !helper.ValidateParam(fp.OrderBy)) || (len(fp.OrderType) > 0 && !helper.ValidateParam(fp.OrderType)) {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	listProduct := []ProductsResponse{}
	listProductRes, totalData, errCode, err := prd.prodRepo.GetListProducts(ctx, repo_products.ProductsParamsTemp{
		Limit:     fp.Limit,
		Page:      fp.Page,
		Name:      fp.Name,
		Posisi:    fp.Posisi,
		UkuranBan: fp.UkuranBan,
		MerkMotor: fp.MerkMotor,
		MerkBan:   fp.MerkBan,
		MinPrice:  fp.MinPrice,
		MaxPrice:  fp.MaxPrice,
		OrderBy:   fp.OrderBy,
		OrderType: fp.OrderType,
	}, custId)

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	for _, val := range listProductRes {
		listProduct = append(listProduct, ProductsResponse{
			KodePLU:        val.KodePLU,
			NamaBarang:     val.NamaBarang,
			Disc:           val.Disc,
			NamaUkuran:     val.NamaUkuran,
			HargaJualFinal: val.HargaJualFinal,
			IsWishList:     val.IsWishlist,
			JenisBan:       val.JenisBan,
			DisplayImage:   prd.baseAssetUrl + cn.ProductDir + val.DisplayImage,
		})
	}

	response.Yay(w, r, ListProductsResponse{
		Products: listProduct,
		DataInfo: DataInfo{
			CurrentPage: fp.Page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(fp.Limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       fp.Limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)
}

func (prd *ProductsHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	var (
		ctx              = r.Context()
		listProductImage = []ProductImage{}
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	productId, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	if productId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidProductID), crashy.ErrInvalidProductID, crashy.Message(crashy.ErrInvalidProductID)), http.StatusBadRequest)
		return
	}

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	product, errCode, err := prd.prodRepo.GetProductDetail(ctx, productId, custId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	prodImg, errCode, err := prd.prodRepo.GetProductImage(ctx, product.KodeBarang)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, val := range prodImg {
		listProductImage = append(listProductImage, ProductImage{
			Url:       prd.baseAssetUrl + cn.ProductDir + val.Url,
			IsDisplay: val.IsDisplay,
		})
	}
	//GET Image list

	response.Yay(w, r, ProductDetailResponse{
		KodePLU:        product.KodePLU,
		NamaBarang:     product.NamaBarang,
		Disc:           product.Disc,
		NamaUkuran:     product.NamaUkuran,
		HargaJualFinal: product.HargaJualFinal,
		IsWishList:     product.IsWishlist,
		JenisBan:       product.JenisBan,
		Posisi:         product.NamaPosisi,
		JenisMotor:     "Bebek",
		TotalTerjual:   0,
		Deskripsi:      product.Deskripsi,
		ImageList:      listProductImage,
		ReviewList: []ProductReview{
			{
				Name:    "Forger",
				Avatar:  prd.baseAssetUrl + cn.UserDir + "profile.png",
				Date:    "2022-05-30",
				Rating:  5,
				Comment: "Barangnya oke banget",
			},
			{
				Name:    "Komang",
				Date:    "2022-05-27",
				Rating:  4,
				Comment: "Enak banget buat sunmorian sambil bawa laptop di tas",
			},
			{
				Name:    "Mr J",
				Date:    "2022-05-30",
				Rating:  5,
				Comment: "these all only dummy comment",
			},
		},
		Kompatibilitas: []MotorCycleCompatibility{
			{
				MerkMotor:    "Honda Vario 125x",
				DisplayImage: prd.baseAssetUrl + cn.MotorBrandDir + "vario125x.png",
			},
			{
				MerkMotor:    "Yamaha Nmax ABS",
				DisplayImage: prd.baseAssetUrl + cn.MotorBrandDir + "nmaxabs.png",
			},
			{
				MerkMotor:    "Honda Beat 125",
				DisplayImage: prd.baseAssetUrl + cn.MotorBrandDir + "hondabeat125.png",
			},
		},
	}, http.StatusOK)
}

func (prd *ProductsHandler) WishlistAdd(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		p        WishlistRequest
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	errCode, err = prd.prodRepo.WishlistAdd(ctx, custId, p.KodePLU)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (prd *ProductsHandler) WishlistRemove(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		p        WishlistRequest
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	errCode, err = prd.prodRepo.WishlistRemove(ctx, custId, p.KodePLU)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) WishlistMe(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		fp       = NewProductsParams(r)
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	listProduct := []ProductsResponse{}
	listProductRes, totalData, errCode, err := prd.prodRepo.WishlistMe(ctx, custId, repo_products.ProductsParamsTemp{
		Limit: fp.Limit,
		Page:  fp.Page,
	})

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	for _, val := range listProductRes {
		listProduct = append(listProduct, ProductsResponse{
			KodePLU:        val.KodePLU,
			NamaBarang:     val.NamaBarang,
			Disc:           val.Disc,
			NamaUkuran:     val.NamaUkuran,
			HargaJualFinal: val.HargaJualFinal,
			IsWishList:     true,
			JenisBan:       val.JenisBan,
			DisplayImage:   prd.baseAssetUrl + cn.ProductDir + val.DisplayImage,
		})
	}

	response.Yay(w, r, ListProductsResponse{
		Products: listProduct,
		DataInfo: DataInfo{
			CurrentPage: fp.Page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(fp.Limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       fp.Limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)

}

func (prd *ProductsHandler) CartAdd(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		p        WishlistRequest
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	cartId, errCode, err := prd.prodRepo.CartCheck(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if cartId == 0 {
		cartId, errCode, err = prd.prodRepo.CartAdd(ctx, authData.Uid)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	}

	cartItemId, qty, errCode, err := prd.prodRepo.CartItemCheck(ctx, cartId, p.KodePLU)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if cartItemId == 0 {
		errCode, err = prd.prodRepo.CartItemAdd(ctx, cartId, p.KodePLU)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	} else {
		errCode, err = prd.prodRepo.CartItemUpdate(ctx, cartItemId, (qty + 1), true)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) CartRemove(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   CartItemRemoveRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := prd.prodRepo.CartItemRemove(ctx, p.CartItemId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) CartUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   CartItemUpdateRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := prd.prodRepo.CartItemUpdate(ctx, p.CartItemId, p.Qty, p.IsSelected)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) CartSelectDeselectAll(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   CartSelectAllRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := prd.prodRepo.CartSelectDeselectAll(ctx, p.CartId, p.IsSelectAll)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) CartMe(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		fp       = NewProductsParams(r)
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if custId == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCode(crashy.ErrInvalidToken), crashy.Message(crashy.ErrInvalidToken)), http.StatusUnauthorized)
		return
	}

	cartId, errCode, err := prd.prodRepo.CartCheck(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if cartId == 0 {
		cartId, errCode, err = prd.prodRepo.CartAdd(ctx, authData.Uid)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	}

	listCartItem := []CartResponse{}
	listProductRes, totalData, errCode, err := prd.prodRepo.CartMe(ctx, cartId, repo_products.ProductsParamsTemp{
		Limit: fp.Limit,
		Page:  fp.Page,
	})

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	for _, val := range listProductRes {
		listCartItem = append(listCartItem, CartResponse{
			CartItemId: val.CartItemId,
			CartItemQty: val.CartItemQty,
			CartItemIsSelected: val.CartItemIsSelected,
			KodePLU:        val.KodePLU,
			NamaBarang:     val.NamaBarang,
			Disc:           val.Disc,
			NamaUkuran:     val.NamaUkuran,
			HargaJualFinal: val.HargaJualFinal,
			JenisBan:       val.JenisBan,
			DisplayImage:   prd.baseAssetUrl + cn.ProductDir + val.DisplayImage,
		})
	}

	response.Yay(w, r, ListItemCartResponse{
		CartId: cartId,
		CartsItem: listCartItem,
		DataInfo: DataInfo{
			CurrentPage: fp.Page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(fp.Limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       fp.Limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)

}
