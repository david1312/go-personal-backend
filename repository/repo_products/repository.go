package repo_products

import (
	"context"
)

type ProductsRepository interface {
	//Products
	GetListProducts(ctx context.Context, fp ProductsParamsTemp, custId int) (res []Products, totalData int, errCode string, err error)
	GetProductDetail(ctx context.Context, id, custId int) (res Products, errCode string, err error)
	GetProductImage(ctx context.Context, productCode string) (res []ProductImage, errCode string, err error)
	GetCustomerId(ctx context.Context, uid string) (custId int, errCode string, err error)
	// Wishlists
	WishlistAdd(ctx context.Context, custId, productId int) (errCode string, err error)
	WishlistRemove(ctx context.Context, custId, productId int) (errCode string, err error)
	WishlistMe(ctx context.Context, custId int, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error)
	// Carts
	CartCheck(ctx context.Context, custUid string) (cartId int, errCode string, err error)
	CartAdd(ctx context.Context, custUid string) (cartId int, errCode string, err error)
	CartSelectDeselectAll(ctx context.Context, cartId int, isSelectAll bool) (errCode string, err error)
	CartItemCheck(ctx context.Context, cartId, productId int) (cartItemId, qty int, errCode string, err error)
	CartItemAdd(ctx context.Context, cartId, productId int) (errCode string, err error)
	CartItemUpdate(ctx context.Context, cartItemId, qty int, isSelected bool) (errCode string, err error)
	CartItemRemove(ctx context.Context, cartItemId int) (errCode string, err error)
	CartMe(ctx context.Context, cartId int, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error)

	//merchant
	DeleteProductById(ctx context.Context, productId int) (errCode string, err error)
	AddProduct(ctx context.Context, sku, name, brandId, tireType, size, price, stock, description string, photoList []string) (errCode string, err error)
}
