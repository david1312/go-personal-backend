package repo_products

import (
	"context"
)

type ProductsRepository interface {
	GetListProducts(ctx context.Context, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error)
}
