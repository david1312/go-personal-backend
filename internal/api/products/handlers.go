package products

import (
	"errors"
	"fmt"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/constants"
	cn "semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/repository/repo_master_data"
	"semesta-ban/repository/repo_products"
	"strconv"
	"strings"

	localMdl "semesta-ban/internal/api/middleware"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type ProductsHandler struct {
	db           *sqlx.DB
	prodRepo     repo_products.ProductsRepository
	mdRepo       repo_master_data.MasterDataRepository
	baseAssetUrl string
	uploadPath   string
	imgMaxSize   int
}

func NewProductsHandler(db *sqlx.DB, pr repo_products.ProductsRepository, md repo_master_data.MasterDataRepository, baseAssetUrl, uploadPath string, maxSize int) *ProductsHandler {
	return &ProductsHandler{db: db, prodRepo: pr, baseAssetUrl: baseAssetUrl, uploadPath: uploadPath, mdRepo: md, imgMaxSize: maxSize}
}

func (prd *ProductsHandler) GetListProducts(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		fp        GetProductsRequest
		authData  = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		arrUkuran = []string{}
	)

	if len(fp.OrderBy) > 0 && !helper.ValidateParam(fp.OrderBy) {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	if err := render.Bind(r, &fp); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	limit := fp.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	page := fp.Page
	if page < 1 {
		page = 1
	}

	custId, errCode, err := prd.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if len(fp.MerkMotor) > 0 && fp.IdMotor == 0 {
		dataListUkuran, errCode, err := prd.mdRepo.GetListUkuranBanByBrandMotor(ctx, fp.MerkMotor)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
		for _, v := range dataListUkuran {
			arrUkuran = append(arrUkuran, v.Id)
		}
	}

	if fp.IdMotor > 0 {
		dataListUkuran, errCode, err := prd.mdRepo.GetListUkuranBanByMotor(ctx, fp.IdMotor)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
		for _, v := range dataListUkuran {
			arrUkuran = append(arrUkuran, v.Id)
		}
	}

	listProduct := []ProductsResponse{}
	listProductRes, totalData, errCode, err := prd.prodRepo.GetListProducts(ctx, repo_products.ProductsParamsTemp{
		Limit:     limit,
		Page:      page,
		Name:      fp.Name,
		UkuranBan: fp.UkuranBan,
		MerkBan:   fp.MerkBan,
		MinPrice:  fp.MinPrice,
		MaxPrice:  fp.MaxPrice,
		OrderBy:   fp.OrderBy,
		ArrUkuran: arrUkuran,
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
			IsWishList:     false,
			JenisBan:       val.JenisBan,
			DisplayImage:   prd.baseAssetUrl + cn.ProductDir + val.DisplayImage,
			IdTireBrand:    val.IDMerk,
			Stock:          val.StockAll,
			Deskripsi:      val.Deskripsi,
		})
	}

	response.Yay(w, r, ListProductsResponse{
		Products: listProduct,
		DataInfo: DataInfo{
			CurrentPage: page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)
}

func (prd *ProductsHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	var (
		ctx              = r.Context()
		listProductImage = []ProductImage{}
		authData         = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		resCompatibilty  = []MotorCycleCompatibility{}
		arrCompatibilty  = []string{}
		topComment       = []ProductReview{}
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
			Id:        val.Id,
			Url:       prd.baseAssetUrl + cn.ProductDir + val.Url,
			IsDisplay: val.IsDisplay, //todo fix tipe data
		})
	}
	//GET Image list

	//get compatibility
	compatibleList, errCode, err := prd.prodRepo.GetProductCompatibility(ctx, product.NamaUkuran)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	for _, v := range compatibleList {
		resCompatibilty = append(resCompatibilty, MotorCycleCompatibility{
			MerkMotor:    v.Motor,
			DisplayImage: prd.baseAssetUrl + cn.MotorCategoryDir + v.DisplayImage,
		})
		arrCompatibilty = append(arrCompatibilty, v.Motor)
	}

	//get top comment
	commentList, errCode, err := prd.prodRepo.GetTopCommentOutlet(ctx) //TODO IMPORTANT IMPROVE LATER WITH ID OUTLET
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, v := range commentList {
		avatar := ""

		if len(v.Avatar) > 0 && v.Avatar[:3] == "pic" {
			avatar = prd.baseAssetUrl + cn.UserDir + v.Avatar
		} else if len(v.Avatar) > 0 && v.Avatar[:3] != "pic" {
			avatar = v.Avatar
		}
		topComment = append(topComment, ProductReview{
			Name:    v.Name,
			Avatar:  avatar,
			Date:    v.Date,
			Rating:  v.Rating,
			Comment: v.Comment,
		})
	}

	response.Yay(w, r, ProductDetailResponse{
		KodePLU:        product.KodePLU,
		NamaBarang:     product.NamaBarang,
		Disc:           product.Disc,
		NamaUkuran:     product.NamaUkuran,
		HargaJualFinal: product.HargaJualFinal,
		IsWishList:     product.IsWishlist,
		JenisBan:       product.JenisBan,
		Posisi:         product.NamaPosisi,
		JenisMotor:     strings.Join(arrCompatibilty, ","),
		TotalTerjual:   0,
		Deskripsi:      product.Deskripsi,
		ImageList:      listProductImage,
		ReviewList:     topComment,
		Kompatibilitas: resCompatibilty,
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
			CartItemId:         val.CartItemId,
			CartItemQty:        val.CartItemQty,
			CartItemIsSelected: val.CartItemIsSelected,
			KodePLU:            val.KodePLU,
			NamaBarang:         val.NamaBarang,
			Disc:               val.Disc,
			NamaUkuran:         val.NamaUkuran,
			HargaJualFinal:     val.HargaJualFinal,
			JenisBan:           val.JenisBan,
			DisplayImage:       prd.baseAssetUrl + cn.ProductDir + val.DisplayImage,
		})
	}

	response.Yay(w, r, ListItemCartResponse{
		CartId:    cartId,
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

func (prd *ProductsHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   ProductCommonRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := prd.prodRepo.DeleteProductById(ctx, p.Id)

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) AddProduct(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = r.Context()
		sku          = r.FormValue("sku")
		name         = r.FormValue("name")
		brandId      = r.FormValue("brand_id")
		tireType     = r.FormValue("tire_type")
		size         = r.FormValue("size")
		price        = r.FormValue("price")
		stock        = r.FormValue("stock")
		description  = r.FormValue("description")
		fileNameList = []string{}
	)

	// validate input
	if len(sku) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "sku cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(name) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "name cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(brandId) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "brand cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(tireType) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "tire type cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(size) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "size cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(price) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "price cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(stock) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "price cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(description) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "description cannot be blank"), http.StatusBadRequest)
		return
	}

	//check all file size before uploading
	for _, fh := range r.MultipartForm.File["photos"] {
		if fh.Size > int64(helper.ConvertFileSizeToMb(prd.imgMaxSize)) {
			errMsg := fmt.Sprintf("%s%v mb", crashy.Message(crashy.ErrCode(crashy.ErrExceededFileSize)), prd.imgMaxSize)
			response.Nay(w, r, crashy.New(errors.New(crashy.ErrExceededFileSize), crashy.ErrExceededFileSize, errMsg), http.StatusBadRequest)
			return
		}

	}

	fileNameList, errCode, err := helper.UploadImage(r, "photos", prd.uploadPath, constants.ProductDir)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = prd.prodRepo.AddProduct(ctx, sku, name, brandId, tireType, size, price, stock, description, fileNameList)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) EPProductUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   UpdateProductRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := prd.prodRepo.ProductUpdate(ctx, repo_products.UpdateProductParam{
		Id:          p.Id,
		Name:        p.Name,
		IdTIreBrand: p.IdTIreBrand,
		TireType:    p.TireType,
		Size:        p.Size,
		Price:       p.Price,
		Stock:       p.Stock,
		Description: p.Description,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) EPProductAddImage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = r.Context()
		id           = r.FormValue("id")
		fileNameList = []string{}
	)

	// validate input
	if len(id) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "id cannot be blank"), http.StatusBadRequest)
		return
	}

	//get data product
	product, errCode, err := prd.prodRepo.ProductDetailMerchant(ctx, id)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//check all file size before uploading
	for _, fh := range r.MultipartForm.File["photos"] {
		if fh.Size > int64(helper.ConvertFileSizeToMb(prd.imgMaxSize)) {
			errMsg := fmt.Sprintf("%s%v mb", crashy.Message(crashy.ErrCode(crashy.ErrExceededFileSize)), prd.imgMaxSize)
			response.Nay(w, r, crashy.New(errors.New(crashy.ErrExceededFileSize), crashy.ErrExceededFileSize, errMsg), http.StatusBadRequest)
			return
		}

	}

	fileNameList, errCode, err = helper.UploadImage(r, "photos", prd.uploadPath, constants.ProductDir)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = prd.prodRepo.ProductAddImage(ctx, product.KodeBarang, fileNameList)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) EpProductDeleteImage(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   DeleteProductImageReq
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	imgDetail, errCode, err := prd.prodRepo.ProductDetailImage(ctx, p.Id)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if imgDetail.Count == 1 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrMinimumPhoto), crashy.ErrMinimumPhoto, crashy.Message(crashy.ErrMinimumPhoto)), http.StatusBadRequest)
		return
	}

	errCode, err = prd.prodRepo.ProductRemoveImage(ctx, p.Id, imgDetail.KodeBarang, imgDetail.Url, prd.uploadPath, cn.ProductDir, imgDetail.IsDisplayFixed)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (prd *ProductsHandler) EpProductUpdateImage(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = r.FormValue("id_image")
	)

	// validate input
	if len(id) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "id image cannot be blank"), http.StatusBadRequest)
		return
	}

	idInt, _ := strconv.Atoi(id)

	imgDetail, errCode, err := prd.prodRepo.ProductDetailImage(ctx, idInt)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	fileName, errCode, err := helper.UploadSingleImage(r, "photo", prd.uploadPath, cn.ProductDir, prd.imgMaxSize)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = prd.prodRepo.ProductUpdateImage(ctx, idInt, fileName, imgDetail.Url, prd.uploadPath, cn.ProductDir)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}
