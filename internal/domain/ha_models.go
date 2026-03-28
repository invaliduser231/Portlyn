package domain

import "time"

type StoredTLSCertificate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Domain       string    `gorm:"uniqueIndex;size:255;not null" json:"domain"`
	IssuerKey    string    `gorm:"size:255;index;not null" json:"issuer_key"`
	Certificate  []byte    `gorm:"not null" json:"-"`
	PrivateKey   []byte    `gorm:"not null" json:"-"`
	MetadataJSON string    `gorm:"type:text" json:"metadata_json"`
	NotBefore    time.Time `gorm:"index;not null" json:"not_before"`
	NotAfter     time.Time `gorm:"index;not null" json:"not_after"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DistributedKV struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Bucket     string    `gorm:"size:64;index:idx_distributed_kv_bucket_key,unique;not null" json:"bucket"`
	Key        string    `gorm:"size:1024;index:idx_distributed_kv_bucket_key,unique;not null" json:"key"`
	Value      []byte    `gorm:"not null" json:"-"`
	ModifiedAt time.Time `gorm:"index;not null" json:"modified_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DistributedLock struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:512;not null" json:"name"`
	Owner     string    `gorm:"size:128;index;not null" json:"owner"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
