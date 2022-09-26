package master_data

import (
	"errors"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/repository/repo_master_data"
	"sort"
	"strconv"
	"strings"

	cn "semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type MasterDataHandler struct {
	db           *sqlx.DB
	mdRepo       repo_master_data.MasterDataRepository
	baseAssetUrl string
	uploadPath   string
	imgMaxSize   int
}

func NewMasterDataHandler(db *sqlx.DB, md repo_master_data.MasterDataRepository, baseAssetUrl, uploadPth string, imgMaxSize int) *MasterDataHandler {
	return &MasterDataHandler{db: db, mdRepo: md, baseAssetUrl: baseAssetUrl, uploadPath: uploadPth, imgMaxSize: imgMaxSize}
}

func (md *MasterDataHandler) GetListMerkBan(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		listMerkBan = []MerkBan{}
	)

	data, errCode, err := md.mdRepo.GetListMerkBan(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, val := range data {
		listMerkBan = append(listMerkBan, MerkBan{
			IdMerk: val.IdMerk,
			Merk:   val.Merk,
			Icon:   md.baseAssetUrl + cn.TireBrandDir + val.Icon,
		})
	}

	response.Yay(w, r, listMerkBan, http.StatusOK)
}

func (md *MasterDataHandler) GetListOutlet(w http.ResponseWriter, r *http.Request) {

	response.Yay(w, r, []Outlet{
		{
			Id:        1,
			Name:      "Semesta Ban",
			Address:   "Jl Raya Kuningan Losari KM 39,5, Desa, Bojongnegara, Kabupaten Cirebon, Jawa Barat 45188",
			Latitude:  -6.8909125,
			Longitude: 108.7525081,
			MapUrl:    "https://goo.gl/maps/e2HJnDKKfzuCeMqZ9",
		},
	}, http.StatusOK)
}

func (md *MasterDataHandler) GetListGender(w http.ResponseWriter, r *http.Request) {

	response.Yay(w, r, []Gender{
		{
			Value: cn.Male,
		},
		{
			Value: cn.Female,
		},
		{
			Value: cn.OtherGender,
		},
	}, http.StatusOK)
}

func (md *MasterDataHandler) GetListSortBy(w http.ResponseWriter, r *http.Request) {
	response.Yay(w, r, ListSortBy, http.StatusOK)
}

func (md *MasterDataHandler) GetListSizeBan(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		listSizeBan = []ListUkuranBan{}
	)

	data, errCode, err := md.mdRepo.GetListUkuranBan(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	mappedData := make(map[int][]UkuranBanTemp)
	for _, m := range data {
		mappedData[m.Ranking] = append(mappedData[m.Ranking], UkuranBanTemp{RingBan: m.UkuranRing, Ukuran: m.Id})
	}
	keys := make([]int, 0, len(mappedData))
	for k := range mappedData {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		tempListSize := []UkuranBan{}
		tempIdRing := ""
		for _, val := range mappedData[k] {
			tempListSize = append(tempListSize, UkuranBan{
				Ukuran: val.Ukuran,
			})
			tempIdRing = val.RingBan
		}

		listSizeBan = append(listSizeBan, ListUkuranBan{
			RingBan:    tempIdRing,
			ListUkuran: tempListSize,
		},
		)
	}

	response.Yay(w, r, listSizeBan, http.StatusOK)

}

func (md *MasterDataHandler) GetListMerkMotor(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		listMerkMotor = []MerkMotor{}
	)

	data, errCode, err := md.mdRepo.GetListMerkMotor(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, val := range data {
		listMerkMotor = append(listMerkMotor, MerkMotor{
			Id:   val.Id,
			Nama: val.Nama,
			Icon: md.baseAssetUrl + cn.MotorBrandDir + val.Icon,
		})
	}

	response.Yay(w, r, listMerkMotor, http.StatusOK)
}

func (md *MasterDataHandler) GetListMotorByBrand(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		listMotor = []ListMotor{}
	)
	idBrandMotor, err := strconv.Atoi(r.URL.Query().Get("id_brand_motor"))

	if err != nil {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	if idBrandMotor == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	data, errCode, err := md.mdRepo.GetListMotorByBrand(ctx, idBrandMotor)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	mappedData := make(map[string][]Motor)
	for _, m := range data {
		mappedData[m.CategoryName] = append(mappedData[m.CategoryName], Motor{
			Id:   m.Id,
			Nama: m.Name,
			Icon: md.baseAssetUrl + cn.MotorDir + m.Icon,
		})
	}
	keys := make([]string, 0, len(mappedData))
	for k := range mappedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tempListMotor := []Motor{}
		for _, val := range mappedData[k] {
			tempListMotor = append(tempListMotor, Motor{
				Id:   val.Id,
				Nama: val.Nama,
				Icon: val.Icon,
			})
		}

		listMotor = append(listMotor, ListMotor{
			Category:  k,
			ListMotor: tempListMotor,
		},
		)
	}

	response.Yay(w, r, listMotor, http.StatusOK)
}

func (md *MasterDataHandler) GetListPaymentMethod(w http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		listPaymentMethod = []ListPaymentMethod{}
	)

	data, errCode, err := md.mdRepo.GetListPaymentMethod(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	mappedData := make(map[string][]PaymentMethod)
	for _, m := range data {
		mappedData[m.CategoryName] = append(mappedData[m.CategoryName], PaymentMethod{
			Id:          m.Id,
			Description: m.Description,
			IsDefault:   m.IsDefault,
			Icon:        md.baseAssetUrl + cn.PaymentMethod + m.Icon,
		})
	}
	keys := make([]string, 0, len(mappedData))
	for k := range mappedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tempListPayment := []PaymentMethod{}
		for _, val := range mappedData[k] {
			tempListPayment = append(tempListPayment, PaymentMethod{
				Id:          val.Id,
				Description: val.Description,
				IsDefault:   val.IsDefault,
				Icon:        val.Icon,
			})
		}

		listPaymentMethod = append(listPaymentMethod, ListPaymentMethod{
			Category:          k,
			ListPaymentMethod: tempListPayment,
		},
		)
	}

	response.Yay(w, r, listPaymentMethod, http.StatusOK)

}

func (md *MasterDataHandler) GetTopRankMotor(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		listMotor = []Motor{}
	)

	data, errCode, err := md.mdRepo.GetListTopRankpMotor(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, val := range data {
		listMotor = append(listMotor, Motor{
			Id:   val.Id,
			Nama: val.Name,
			Icon: md.baseAssetUrl + cn.MotorDir + val.Icon,
		})
	}

	response.Yay(w, r, listMotor, http.StatusOK)
}

func (md *MasterDataHandler) GetImgAsset(w http.ResponseWriter, r *http.Request) {

	response.Yay(w, r, ImageAssetResponse{
		PromoBannerData: PromoBanner{
			Alt:      "Promo Banner",
			ImageUrl: md.baseAssetUrl + cn.PromoDir + cn.StaticPromoBanner,
		},
	}, http.StatusOK)
}

func (md *MasterDataHandler) GetTireType(w http.ResponseWriter, r *http.Request) {

	response.Yay(w, r, []TireType{
		{Value: "TUBE TYPE"},
		{Value: "TUBE LESS"},
	}, http.StatusOK)
}

func (md *MasterDataHandler) EPAddBrandMotor(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		name = r.FormValue("name")
	)

	// validate input
	if len(name) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "name cannot be blank"), http.StatusBadRequest)
		return
	}

	fileName, errCode, err := helper.UploadSingleImage(r, "icon", md.uploadPath, cn.MotorBrandDir, md.imgMaxSize)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.AddBrandMotor(ctx, name, fileName)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (md *MasterDataHandler) EPRemoveBrandMotor(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   MasterDataCommonRequest
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	brandExists, errCode, err := md.mdRepo.CheckBrandMotorUsed(ctx, p.Id)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if brandExists {
		response.Nay(w, r, crashy.New(err, crashy.ErrBrandMotorUsed, crashy.Message(crashy.ErrBrandMotorUsed)), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.RemoveBrandMotor(ctx, p.Id, md.uploadPath, cn.MotorBrandDir)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (md *MasterDataHandler) EPUpdateBrandMotor(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   UpdateBrandMotorReq
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := md.mdRepo.UpdateBrandMotor(ctx, p.Id, p.Name)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (md *MasterDataHandler) EPUpdateBrandMotorIcon(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = r.FormValue("id")
	)

	// validate input
	if len(id) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "id cannot be blank"), http.StatusBadRequest)
		return
	}
	idMotor, _ := strconv.Atoi(id)
	brandExists, errCode, err := md.mdRepo.CheckBrandMotorExist(ctx, idMotor)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if !brandExists {
		response.Nay(w, r, crashy.New(err, crashy.ErrInvalidBrandMotor, crashy.Message(crashy.ErrInvalidBrandMotor)), http.StatusBadRequest)
		return
	}

	fileName, errCode, err := helper.UploadSingleImage(r, "icon", md.uploadPath, cn.MotorBrandDir, md.imgMaxSize)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.UpdateBrandMotorImage(ctx, idMotor, fileName, md.uploadPath, cn.MotorBrandDir)

	response.Yay(w, r, "success", http.StatusOK)

}

func (md *MasterDataHandler) EPAddTireBrand(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		id      = r.FormValue("id")
		name    = r.FormValue("name")
		ranking = r.FormValue("ranking")
	)

	// validate input
	if len(id) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "name cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(name) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "name cannot be blank"), http.StatusBadRequest)
		return
	}
	if len(ranking) == 0 {
		ranking = "99"
	}

	fileName, errCode, err := helper.UploadSingleImage(r, "icon", md.uploadPath, cn.TireBrandDir, md.imgMaxSize)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.AddTireBrand(ctx, strings.ToUpper(id), name, fileName, ranking)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (md *MasterDataHandler) EPRemoveTireBrand(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   MasterDataCommonRequestSec
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	brandExists, errCode, err := md.mdRepo.CheckTireBrandUsed(ctx, p.Id)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if brandExists {
		response.Nay(w, r, crashy.New(err, crashy.ErrTireBrandUsed, crashy.Message(crashy.ErrTireBrandUsed)), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.RemoveTireBrand(ctx, p.Id, md.uploadPath, cn.TireBrandDir)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}

func (md *MasterDataHandler) EPUpdateTireBrand(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		p   UpdateTireBrandReq
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := md.mdRepo.UpdateTireBrand(ctx, p.Id, p.Name, p.Ranking)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	response.Yay(w, r, "success", http.StatusOK)
}


func (md *MasterDataHandler) EPUpdateTireBrandIcon(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = r.FormValue("id")
	)

	// validate input
	if len(id) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, "id cannot be blank"), http.StatusBadRequest)
		return
	}
	brandExists, errCode, err := md.mdRepo.CheckTireBrandExist(ctx, id)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if !brandExists {
		response.Nay(w, r, crashy.New(err, crashy.ErrInvalidTireBrand, crashy.Message(crashy.ErrInvalidTireBrand)), http.StatusBadRequest)
		return
	}

	fileName, errCode, err := helper.UploadSingleImage(r, "icon", md.uploadPath, cn.TireBrandDir, md.imgMaxSize)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	errCode, err = md.mdRepo.UpdateTireBrandImage(ctx, id, fileName, md.uploadPath, cn.TireBrandDir)

	response.Yay(w, r, "success", http.StatusOK)

}