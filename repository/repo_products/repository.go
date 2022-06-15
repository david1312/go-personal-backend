package repo_products

import (
	"context"
)

type ProductsRepository interface {
	GetListProducts(ctx context.Context, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error)
	GetProductDetail(ctx context.Context, id int) (res Products, errCode string, err error)
	GetProductImage(ctx context.Context, productCode string) (res []ProductImage, errCode string, err error)
	GetCustomerId(ctx context.Context, uid string) (custId int, errCode string, err error)
	WishlistAdd(ctx context.Context, custId, productId int) (errCode string, err error)
	WishlistRemove(ctx context.Context, custId, productId int) (errCode string, err error)
	WishlistMe(ctx context.Context, custId int, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error)
}
