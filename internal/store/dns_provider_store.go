package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type DNSProviderStore struct {
	db *gorm.DB
}

func NewDNSProviderStore(db *gorm.DB) *DNSProviderStore {
	return &DNSProviderStore{db: db}
}

func (s *DNSProviderStore) List(ctx context.Context) ([]domain.DNSProvider, error) {
	var items []domain.DNSProvider
	err := s.db.WithContext(ctx).Order("name asc").Find(&items).Error
	return items, err
}

func (s *DNSProviderStore) Count(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.DNSProvider{}).Count(&count).Error
	return count, err
}

func (s *DNSProviderStore) Create(ctx context.Context, item *domain.DNSProvider) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *DNSProviderStore) GetByID(ctx context.Context, id uint) (*domain.DNSProvider, error) {
	var item domain.DNSProvider
	err := s.db.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *DNSProviderStore) Update(ctx context.Context, item *domain.DNSProvider) error {
	return s.db.WithContext(ctx).Save(item).Error
}

func (s *DNSProviderStore) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&domain.DNSProvider{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
