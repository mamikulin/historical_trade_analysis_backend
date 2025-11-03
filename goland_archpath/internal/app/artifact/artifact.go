package artifact

import "gorm.io/gorm"

type Artifact struct {
    gorm.Model
    Name             string  `gorm:"not null" json:"name"`
    Description      string  `gorm:"type:text" json:"description"`
    IsActive         bool    `gorm:"not null;default:true" json:"is_active"`
    ImageURL         *string `json:"image_url"`
    ProductionCenter string  `gorm:"not null;default:'Unknown'" json:"production_center"`
}
