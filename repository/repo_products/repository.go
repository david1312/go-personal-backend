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
	GetProductCompatibility(ctx context.Context, sizeBan string) (res []MotorCycleCompatibility, errCode string, err error)
	GetTopCommentOutlet(ctx context.Context) (res []ProductReview, errCode string, err error)
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
	AddProduct(ctx context.Context, sku, name, brandId, tireType, strikePrice, size, price, stock, description string, photoList []string) (errCode string, err error)
	ProductUpdate(ctx context.Context, param UpdateProductParam) (errCode string, err error)
	ProductAddImage(ctx context.Context, sku string, photoList []string) (errCode string, err error)
	ProductDetailMerchant(ctx context.Context, id string) (res Products, errCode string, err error)
	ProductDetailImage(ctx context.Context, idImg int) (res ProductImage, errCode string, err error)
	ProductRemoveImage(ctx context.Context, idImg int, kodeBarang, fileName, uploadPath, dirFile string, isDisplay bool) (errCode string, err error)
	ProductUpdateImage(ctx context.Context, idImg int, fileNameNew, fileNameOld, uploadPath, dirFile string) (errCode string, err error)
}
