package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type ClientStore struct {
	db *gorm.DB
}

func NewClientStore(db *gorm.DB) *ClientStore {
	return &ClientStore{db: db}
}

func (s *ClientStore) List(ctx context.Context) ([]domain.Client, error) {
	var clients []domain.Client
	err := s.db.WithContext(ctx).Order("id asc").Find(&clients).Error
	return clients, err
}

func (s *ClientStore) Create(ctx context.Context, client *domain.Client) error {
	return s.db.WithContext(ctx).Create(client).Error
}

func (s *ClientStore) GetByID(ctx context.Context, id uint) (*domain.Client, error) {
	var client domain.Client
	err := s.db.WithContext(ctx).First(&client, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *ClientStore) Update(ctx context.Context, client *domain.Client) error {
	return s.db.WithContext(ctx).Save(client).Error
}

func (s *ClientStore) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&domain.Client{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
