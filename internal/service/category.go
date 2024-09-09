package service

import (
	"context"
	"io"

	"github.com/Genarodaniel/go-grpc-study/internal/database"
	"github.com/Genarodaniel/go-grpc-study/internal/pb"
	"google.golang.org/grpc"
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

func (c *CategoryService) GetCategory(ctx context.Context, request *pb.CategoryGetByIDRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.FindByID(request.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream grpc.ClientStreamingServer[pb.CreateCategoryRequest, pb.CategoryList]) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err != nil && err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream grpc.BidiStreamingServer[pb.CreateCategoryRequest, pb.Category]) error {
	for {
		category, err := stream.Recv()
		if err != nil && err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        category.Name,
			Description: category.Description,
		}); err != nil {
			return err
		}

	}
}
