package service

import (
	"context"

	"github.com/Genarodaniel/go-grpc-study/internal/database"
	"github.com/Genarodaniel/go-grpc-study/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(request.Name, request.Description)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Name:        category.Name,
		Id:          category.ID,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) ListCategories(context.Context, *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	categoryList := pb.CategoryList{}
	for _, category := range categories {
		categoryList.Categories = append(categoryList.Categories, &pb.Category{
			Id:          category.ID,
			Description: category.Description,
			Name:        category.Name,
		})
	}

	return &categoryList, nil
}
