package repo_master_data

import "context"

type MasterDataRepository interface {
	GetListMerkBan(ctx context.Context) (res []MerkBan, errCode string, err error)
}
