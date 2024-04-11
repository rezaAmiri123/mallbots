package queries

// import (
// 	"context"

// 	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
// )


// type GetProduct struct {
// 	ID string
// }

// type GetProductHandler struct {
// 	mall domain.MallRepository
// }

// func NewGetProductHandler(mall domain.MallRepository) GetProductHandler {
// 	return GetProductHandler{mall: mall}
// }

// func (h GetProductHandler) GetProduct(ctx context.Context, query GetProduct) (*domain.MallProduct, error) {
// 	return h.mall.Find(ctx, query.ID)
// }
