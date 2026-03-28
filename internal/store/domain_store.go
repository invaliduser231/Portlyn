package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type DomainStore struct {
	db *gorm.DB
}

func NewDomainStore(db *gorm.DB) *DomainStore {
	return &DomainStore{db: db}
}

func (s *DomainStore) List(ctx context.Context) ([]domain.Domain, error) {
	var items []domain.Domain
	err := s.db.WithContext(ctx).Order("id asc").Find(&items).Error
	return items, err
}

func (s *DomainStore) Create(ctx context.Context, item *domain.Domain) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *DomainStore) Count(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.Domain{}).Count(&count).Error
	return count, err
}

func (s *DomainStore) GetByID(ctx context.Context, id uint) (*domain.Domain, error) {
	var item domain.Domain
	err := s.db.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *DomainStore) Update(ctx context.Context, item *domain.Domain) error {
	return s.db.WithContext(ctx).Save(item).Error
}

func (s *DomainStore) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&domain.Domain{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
