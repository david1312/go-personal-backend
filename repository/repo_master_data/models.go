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

type MotorMD struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	BrandMotor    string `json:"brand_motor"`
	CategoryMotor string `json:"category_motor"`
	Icon          string `json:"icon"`
}

type ListMotorRequestRepo struct {
	Limit           int    `json:"limit"`
	Page            int    `json:"page"`
	Name            string `json:"name"`
	IdBrandMotor    int    `json:"id_brand_motor"`
	IdCategoryMotor int    `json:"id_jenis_motor"`
}

type CategoryMotor struct {
	Id   int    `json:"limit"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
