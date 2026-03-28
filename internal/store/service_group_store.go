package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type ServiceGroupStore struct {
	db *gorm.DB
}

func NewServiceGroupStore(db *gorm.DB) *ServiceGroupStore {
	return &ServiceGroupStore{db: db}
}

func (s *ServiceGroupStore) List(ctx context.Context) ([]domain.ServiceGroup, error) {
	var items []domain.ServiceGroup
	err := s.db.WithContext(ctx).Preload("Services").Preload("Services.Domain").Order("name asc").Find(&items).Error
	return items, err
}

func (s *ServiceGroupStore) Count(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.ServiceGroup{}).Count(&count).Error
	return count, err
}

func (s *ServiceGroupStore) GetByID(ctx context.Context, id uint) (*domain.ServiceGroup, error) {
	var item domain.ServiceGroup
	err := s.db.WithContext(ctx).Preload("Services").Preload("Services.Domain").First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ServiceGroupStore) Create(ctx context.Context, item *domain.ServiceGroup) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *ServiceGroupStore) Update(ctx context.Context, item *domain.ServiceGroup) error {
	return s.db.WithContext(ctx).Omit("Services.*").Save(item).Error
}

func (s *ServiceGroupStore) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&domain.ServiceGroup{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *ServiceGroupStore) ReplaceServices(ctx context.Context, serviceGroupID uint, serviceIDs []uint) error {
	item := &domain.ServiceGroup{ID: serviceGroupID}
	services := make([]domain.Service, 0, len(serviceIDs))
	for _, id := range serviceIDs {
		services = append(services, domain.Service{ID: id})
	}
	return s.db.WithContext(ctx).Model(item).Association("Services").Replace(services)
}

func (s *ServiceGroupStore) AddService(ctx context.Context, serviceGroupID, serviceID uint) error {
	item := &domain.ServiceGroup{ID: serviceGroupID}
	return s.db.WithContext(ctx).Model(item).Association("Services").Append(&domain.Service{ID: serviceID})
}

func (s *ServiceGroupStore) RemoveService(ctx context.Context, serviceGroupID, serviceID uint) error {
	item := &domain.ServiceGroup{ID: serviceGroupID}
	return s.db.WithContext(ctx).Model(item).Association("Services").Delete(&domain.Service{ID: serviceID})
}
