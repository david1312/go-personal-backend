package repo_merchant

type MerchantData struct {
	MerchantId     int    `json:"merchant_id"`
	Username       string `json:"username"`
	OutletId       int    `json:"outlet_id"`
	OutletName     string `json:"outlet_name"`
	OutletAvatar   string `json:"outlet_avatar"`
	OutletEmail    string `json:"outlet_email"`
	CsNumber       string `json:"outlet_cs_number"`
	OutletAddress  string `json:"outlet_address"`
	OutletDistrict string `json:"outlet_district"`
	OutletCity     string `json:"outlet_city"`
	OutletGmapUrl  string `json:"outlet_gmap_url"`
	Password       string `json:"password"`
}
