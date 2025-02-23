package service

import (
	"context"

	"github.com/stipochka/web_service/internal/repository"

	"github.com/stipochka/web_service/internal/models"
)

type GrpcClient interface {
	GetAllRecords(ctx context.Context) ([]models.Record, error)
	GetRecordById(ctx context.Context, id int) (models.Record, error)
}

type Service struct {
	GrpcClient
}

func NewService(client *repository.GRPCRepository) *Service {
	return &Service{
		GrpcClient: NewGRPCService(client),
	}
}
