package services

import (
	"context"
	"log"

	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/repositories"
)

type AdvertisingService interface {
	CreateAd(ctx context.Context, advertising entities.Advertising) (int, error)
	GetByQueyParams(ctx context.Context, query map[string]string) ([]entities.Advertising, error)
}
type advertisingService struct {
	advRepo repositories.AdvertisingRepository
}

func NewAdvertisingService(advRepo repositories.AdvertisingRepository) AdvertisingService {
	return &advertisingService{advRepo: advRepo}
}

func (ad *advertisingService) CreateAd(ctx context.Context, advertising entities.Advertising) (int, error) {
	id, err := ad.advRepo.Store(ctx, advertising)
	log.Println(err)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ad *advertisingService) GetByQueyParams(ctx context.Context, query map[string]string) ([]entities.Advertising, error) {
	return nil, nil
}
