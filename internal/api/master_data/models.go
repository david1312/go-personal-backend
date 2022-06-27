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

type SortBy struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

//todo update list product sort by
var ListSortBy = []SortBy{{
	Value: "latest",
	Label: "Terbaru",
},
	{
		Value: "max_price",
		Label: "Harga Tertinggi",
	},
	{
		Value: "min_price",
		Label: "Harga Terendah",
	},
}

type UkuranBan struct {
	Ukuran string `json:"ukuran"`
}

type ListUkuranBan struct {
	RingBan    string      `json:"ring_ban"`
	ListUkuran []UkuranBan `json:"list_ukuran"`
}

type MerkMotor struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
	Icon string `json:"icon"`
}

type Motor struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
	Icon string `json:"icon"`
}

type ListMotor struct {
	Category  string  `json:"kategori"`
	ListMotor []Motor `json:"list_motor"`
}
