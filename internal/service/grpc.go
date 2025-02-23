package service

import (
	"context"

	"github.com/stipochka/web_service/internal/repository"

	"github.com/stipochka/web_service/internal/models"
)

type GrpcService struct {
	grpcClient *repository.GRPCRepository
}

func NewGRPCService(cl *repository.GRPCRepository) *GrpcService {
	return &GrpcService{
		grpcClient: cl,
	}
}

//GetAllRecords(ctx context.Context)
//GetRecordById(ctx context.Context, id int)

func (g *GrpcService) GetAllRecords(ctx context.Context) ([]models.Record, error) {
	return g.grpcClient.GetAllRecords(ctx)
}

func (g *GrpcService) GetRecordById(ctx context.Context, id int) (models.Record, error) {
	return g.grpcClient.GetRecordById(ctx, id)
}
