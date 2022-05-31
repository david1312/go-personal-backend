package master_data

import (
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/repository/repo_master_data"

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
