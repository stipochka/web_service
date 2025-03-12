package repository

import (
	"context"

	"github.com/stipochka/web_service/internal/models"
	"google.golang.org/grpc"

	database "github.com/stipochka/protos/gen/go/db"
)

type gRPCclient struct {
	client database.DatabaseClient
}

//GetRecordById(ctx context.Context, id int) (models.Record, error)
//GetAllRecords(ctx context.Context) ([]models.Record, error)

func NewGRPCClient(cl *grpc.ClientConn) *gRPCclient {
	dbClient := database.NewDatabaseClient(cl)
	return &gRPCclient{
		client: dbClient,
	}
}

func (g *gRPCclient) GetRecordById(ctx context.Context, id int) (models.Record, error) {
	record, err := g.client.GetRecordByID(ctx, &database.GetByIdRequest{RecordID: int64(id)})
	if err != nil {
		return models.Record{}, err
	}

	return ToModelsRecord(record), nil
}

func (g *gRPCclient) GetAllRecords(ctx context.Context) ([]models.Record, error) {
	recordsRes := make([]models.Record, 0)

	records, err := g.client.GetAllRecords(ctx, &database.GetAllRecordsRequest{})
	if err != nil {
		return recordsRes, err
	}

	for _, record := range records.Records {
		recordsRes = append(recordsRes, ToModelsRecord(record))
	}

	return recordsRes, nil
}

func ToModelsRecord(rec *database.RecordResponse) models.Record {
	return models.Record{
		ID:            int(rec.GetId()),
		Status:        rec.GetStatus(),
		PollingPeriod: int(rec.GetPollingPeriod()),
		Temperature:   int(rec.GetTemperature()),
		LampStatus:    rec.GetLampStatus(),
		Voltage:       int(rec.GetVoltage()),
		Thresholds: models.Thresholds{
			Temperature: ToModelsSensorData(rec.GetThresholds().GetTemperature()),
			Humidity:    ToModelsSensorData(rec.GetThresholds().GetHumidity()),
			Voltage:     ToModelsSensorData(rec.GetThresholds().GetVoltage()),
		},
	}
}

func ToModelsSensorData(records []*database.SensorData) []models.SensorData {
	res := make([]models.SensorData, 0)
	for _, record := range records {
		res = append(res, models.SensorData{
			Value:         int(record.GetValue()),
			PollingPeriod: int(record.GetPollingPeriod()),
		})
	}
	return res
}
