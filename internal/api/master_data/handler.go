package master_data

import (
	"errors"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/repository/repo_master_data"
	"sort"
	"strconv"

	cn "semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"

	"github.com/jmoiron/sqlx"
)

type MasterDataHandler struct {
	db           *sqlx.DB
	mdRepo       repo_master_data.MasterDataRepository
	baseAssetUrl string
}

func NewMasterDataHandler(db *sqlx.DB, md repo_master_data.MasterDataRepository, baseAssetUrl string) *MasterDataHandler {
	return &MasterDataHandler{db: db, mdRepo: md, baseAssetUrl: baseAssetUrl}
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
			Latitude:  -6.890816046622402,
			Longitude: 108.75284502302158,
			MapUrl:    "https://goo.gl/maps/gEqkgKohXGhU8jXv8",
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

	mappedData := make(map[string][]string)
	for _, m := range data {
		mappedData[m.UkuranRing] = append(mappedData[m.UkuranRing], m.Id)
	}
	keys := make([]string, 0, len(mappedData))
	for k := range mappedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tempListSize := []UkuranBan{}
		for _, val := range mappedData[k] {
			tempListSize = append(tempListSize, UkuranBan{
				Ukuran: val,
			})
		}

		listSizeBan = append(listSizeBan, ListUkuranBan{
			RingBan:    k,
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
			Icon: md.baseAssetUrl + cn.TireBrandDir + val.Icon,
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
			Icon: m.Icon,
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
