package store

import (
	"context"
	"time"

	"gorm.io/gorm"

	"portlyn/internal/domain"
)

type AuditListParams struct {
	Limit        int
	Offset       int
	UserID       *uint
	ResourceType string
	ResourceID   *uint
	ActionLike   string
	RequestID    string
	Method       string
	StatusCode   *int
	Host         string
	From         *time.Time
	To           *time.Time
}

type AuditStore struct {
	db *gorm.DB
}

func NewAuditStore(db *gorm.DB) *AuditStore {
	return &AuditStore{db: db}
}

func (s *AuditStore) Create(ctx context.Context, item *domain.AuditLog) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *AuditStore) CreateBatch(ctx context.Context, items []domain.AuditLog) error {
	if len(items) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Create(&items).Error
}

func (s *AuditStore) List(ctx context.Context, params AuditListParams) ([]domain.AuditLog, int64, error) {
	query := s.db.WithContext(ctx).Model(&domain.AuditLog{})
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.ResourceType != "" {
		query = query.Where("resource_type = ?", params.ResourceType)
	}
	if params.ResourceID != nil {
		query = query.Where("resource_id = ?", *params.ResourceID)
	}
	if params.ActionLike != "" {
		query = query.Where("action LIKE ?", params.ActionLike)
	}
	if params.RequestID != "" {
		query = query.Where("request_id = ?", params.RequestID)
	}
	if params.Method != "" {
		query = query.Where("method = ?", params.Method)
	}
	if params.StatusCode != nil {
		query = query.Where("status_code = ?", *params.StatusCode)
	}
	if params.Host != "" {
		query = query.Where("host = ?", params.Host)
	}
	if params.From != nil {
		query = query.Where("timestamp >= ?", *params.From)
	}
	if params.To != nil {
		query = query.Where("timestamp <= ?", *params.To)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit <= 0 || params.Limit > 200 {
		params.Limit = 50
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	var items []domain.AuditLog
	err := query.Order("timestamp desc").Limit(params.Limit).Offset(params.Offset).Find(&items).Error
	return items, total, err
}

func (s *AuditStore) CountByActionLikeSince(ctx context.Context, actionLike string, since time.Time) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.AuditLog{}).
		Where("action LIKE ? AND timestamp >= ?", actionLike, since).
		Count(&count).Error
	return count, err
}
