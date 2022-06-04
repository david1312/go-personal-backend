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
	)

	if (len(fp.OrderBy) > 0 && !helper.ValidateParam(fp.OrderBy)) || (len(fp.OrderType) > 0 && !helper.ValidateParam(fp.OrderType)) {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
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
			IsWishList:     false,
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

	product, errCode, err := prd.prodRepo.GetProductDetail(ctx, productId)
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
		IsWishList:     false, //todo to be implemented
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
