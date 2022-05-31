package master_data

type Outlet struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	MapUrl    string  `json:"gmap_url"`
}

type Gender struct {
	Value string `json:"value"`
}

type MerkBan struct {
	IdMerk string `json:"id_merk"`
	Merk   string `json:"merk"`
	Icon   string `json:"icon"`
}
