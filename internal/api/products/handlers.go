package products

import (
	"errors"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/repository/repo_products"

	"github.com/jmoiron/sqlx"
)

type ProductsHandler struct {
	db       *sqlx.DB
	prodRepo repo_products.ProductsRepository
}

func NewProductsHandler(db *sqlx.DB, pr repo_products.ProductsRepository) *ProductsHandler {
	return &ProductsHandler{db: db, prodRepo: pr}
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
			HargaJual:      val.HargaJual,
			HargaJualFinal: val.HargaJualFinal,
			IsWishList:     false,
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
