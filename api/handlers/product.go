package handlers

import (
	"context"
	"fmt"

	"github.com/TiagoBehencks/gRPC/api/db"
	"github.com/TiagoBehencks/gRPC/api/pb"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
}

func NewProductServer() *ProductServer {
	return &ProductServer{}
}

func (s *ProductServer) Create(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	var id int32
	err := db.Pool.QueryRow(ctx,
		"INSERT INTO products (name, price, quantity) VALUES ($1, $2, $3) RETURNING id",
		req.Name, req.Price, req.Quantity,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return &pb.Product{
		Id:       fmt.Sprintf("%d", id),
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	}, nil
}

func (s *ProductServer) Get(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	var product pb.Product
	err := db.Pool.QueryRow(ctx,
		"SELECT id, name, price, quantity FROM products WHERE id = $1",
		req.Id,
	).Scan(&product.Id, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	return &product, nil
}

func (s *ProductServer) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	result, err := db.Pool.Exec(ctx,
		"UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4",
		req.Name, req.Price, req.Quantity, req.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("product not found")
	}
	return &pb.Product{
		Id:       req.Id,
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	}, nil
}

func (s *ProductServer) Delete(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	result, err := db.Pool.Exec(ctx, "DELETE FROM products WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("product not found")
	}
	return &pb.Empty{}, nil
}

func (s *ProductServer) List(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	rows, err := db.Pool.Query(ctx, "SELECT id, name, price, quantity FROM products")
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var p pb.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &p)
	}
	return &pb.ListProductsResponse{Products: products}, nil
}