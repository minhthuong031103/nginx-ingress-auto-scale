package server

import (
	"context"
	"productsvc/internal/dal"
	pb "productsvc/internal/generated/product"
	"productsvc/internal/helper"

	"github.com/gocql/gocql"
)

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	productID := gocql.TimeUUID()
	if req.ImageUrl == "" {
		req.ImageUrl = defaultImageURL
	}
	product := dal.Product{
		ProductID:   productID,
		ProductName: req.ProductName,
		Price:       req.Price,
		Description: req.Description,
		Quantity:    req.Quantity,
		Sold:        0,
		ImageURL:    req.ImageUrl,
		CreatedAt:   helper.GetTimeNowGMT7(),
		UpdatedAt:   "",
		DeletedAt:   "",
	}
	err := s.ProductDAL.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return &pb.Product{ProductId: productID.String()}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	productID, err := gocql.ParseUUID(req.ProductId)
	if err != nil {
		return nil, err
	}
	product, err := s.ProductDAL.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	return &pb.Product{
		ProductId:   product.ProductID.String(),
		ProductName: product.ProductName,
		Price:       product.Price,
		Description: product.Description,
		Quantity:    product.Quantity,
		Sold:        product.Sold,
		ImageUrl:    product.ImageURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		DeletedAt:   product.DeletedAt,
	}, nil
}

func (s *Server) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	products, err := s.ProductDAL.GetAllProducts()
	if err != nil {
		return nil, err
	}
	var productResponses []*pb.Product
	for _, product := range products {
		productResponses = append(productResponses, &pb.Product{
			ProductId:   product.ProductID.String(),
			ProductName: product.ProductName,
			Price:       product.Price,
			Description: product.Description,
			Quantity:    product.Quantity,
			Sold:        product.Sold,
			ImageUrl:    product.ImageURL,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			DeletedAt:   product.DeletedAt,
		})
	}
	return &pb.GetAllProductsResponse{
		Products: productResponses,
		Total:    int32(len(productResponses)),
	}, nil
}

func (s *Server) UpdateProductQuantityAndSold(ctx context.Context, req *pb.UpdateProductQuantityAndSoldRequest) (*pb.Product, error) {
	productID, err := gocql.ParseUUID(req.ProductId)
	if err != nil {
		return nil, err
	}
	err = s.ProductDAL.UpdateProductQuantityAndSold(productID, req.Quantity, req.Sold)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := s.ProductDAL.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	return &pb.Product{
		ProductId:   updatedProduct.ProductID.String(),
		ProductName: updatedProduct.ProductName,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		Quantity:    updatedProduct.Quantity,
		Sold:        updatedProduct.Sold,
		ImageUrl:    updatedProduct.ImageURL,
		CreatedAt:   updatedProduct.CreatedAt,
		UpdatedAt:   updatedProduct.UpdatedAt,
		DeletedAt:   updatedProduct.DeletedAt,
	}, nil
}
