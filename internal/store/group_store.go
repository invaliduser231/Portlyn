package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type GroupStore struct {
	db *gorm.DB
}

func NewGroupStore(db *gorm.DB) *GroupStore {
	return &GroupStore{db: db}
}

func (s *GroupStore) List(ctx context.Context) ([]domain.Group, error) {
	var items []domain.Group
	err := s.db.WithContext(ctx).Order("name asc").Find(&items).Error
	return items, err
}

func (s *GroupStore) Count(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.Group{}).Count(&count).Error
	return count, err
}

func (s *GroupStore) GetByID(ctx context.Context, id uint) (*domain.Group, error) {
	var item domain.Group
	err := s.db.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *GroupStore) Create(ctx context.Context, item *domain.Group) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *GroupStore) Update(ctx context.Context, item *domain.Group) error {
	return s.db.WithContext(ctx).Save(item).Error
}

func (s *GroupStore) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&domain.Group{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *GroupStore) CountMembers(ctx context.Context, groupID uint) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.GroupMembership{}).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}

func (s *GroupStore) ListMembers(ctx context.Context, groupID uint) ([]domain.User, error) {
	var users []domain.User
	err := s.db.WithContext(ctx).
		Table("users").
		Joins("join group_memberships on group_memberships.user_id = users.id").
		Where("group_memberships.group_id = ?", groupID).
		Order("users.email asc").
		Find(&users).Error
	return users, err
}

func (s *GroupStore) AddMember(ctx context.Context, groupID, userID uint) error {
	membership := &domain.GroupMembership{GroupID: groupID, UserID: userID}
	return s.db.WithContext(ctx).Where(domain.GroupMembership{GroupID: groupID, UserID: userID}).FirstOrCreate(membership).Error
}

func (s *GroupStore) RemoveMember(ctx context.Context, groupID, userID uint) error {
	result := s.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&domain.GroupMembership{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *GroupStore) ListGroupIDsForUser(ctx context.Context, userID uint) ([]uint, error) {
	var ids []uint
	err := s.db.WithContext(ctx).Model(&domain.GroupMembership{}).Where("user_id = ?", userID).Pluck("group_id", &ids).Error
	return ids, err
}
