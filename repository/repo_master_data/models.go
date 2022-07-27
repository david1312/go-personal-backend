package repo_master_data

type MerkBan struct {
	IdMerk  string `json:"id_merk"`
	Merk    string `json:"merk"`
	Icon    string `json:"icon"`
	Ranking int    `json:"ranking"`
}

type UkuranRingBan struct {
	Id          string `json:"id"`
	IdRingBan   int    `json:"id_ring_ban"`
	IdUkuranBan string `json:"id_ukuran_ban"`
	UkuranRing  string `json:"UkuranRing"`
	Ranking     int    `json:"ranking"`
}

type MerkMotor struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
	Icon string `json:"icon"`
}

type Motor struct {
	Id           int    `json:"id"`
	Name         string `json:"nama"`
	Icon         string `json:"icon"`
	CategoryName string `json:"category_name"`
}

type PaymentMethod struct {
	Id           string `json:"id"`
	Description  string `json:"description"`
	IsDefault    bool   `json:"is_default"`
	Icon         string `json:"icon"`
	CategoryName string `json:"category_name"`
}
