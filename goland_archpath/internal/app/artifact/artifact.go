package artifact

import (
    "time"
)

type Artifact struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`

    Name             string  `gorm:"not null" json:"name"`
    Description      string  `gorm:"type:text" json:"description"`
    IsActive         bool    `gorm:"not null;default:true" json:"is_active"`
    ImageURL         *string `json:"image_url"`
    ProductionCenter string  `gorm:"not null;default:'Unknown'" json:"production_center"`
}
