package artifact

import (
    "time"
    "gorm.io/gorm"
)

type Artifact struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

    Name             string  `gorm:"not null" json:"name"`
    Description      string  `gorm:"type:text" json:"description"`
    ImageURL         *string `json:"image_url"`
    ProductionCenter string  `gorm:"not null;default:'Unknown'" json:"production_center"`
}
