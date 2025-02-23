package repository

import (
	"context"

	_ "github.com/stipochka/protos/gen/go/db"
	"github.com/stipochka/web_service/internal/models"
	"google.golang.org/grpc"
)

type GrpcClient interface {
	GetRecordById(ctx context.Context, id int) (models.Record, error)
	GetAllRecords(ctx context.Context) ([]models.Record, error)
}

type GRPCRepository struct {
	GrpcClient
}

func NewGRPCRepository(cl *grpc.ClientConn) *GRPCRepository {
	return &GRPCRepository{
		GrpcClient: NewGRPCClient(cl),
	}
}
