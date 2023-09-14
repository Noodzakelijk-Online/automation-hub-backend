package model

import (
	"fmt"
	"github.com/google/uuid"
)

type Automation struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name     string    `gorm:"type:varchar(50);uniqueIndex:idx_name_position"`
	Image    string    `gorm:"type:varchar(255)"`
	Host     string    `gorm:"type:varchar(50)"`
	Port     int       `gorm:"check:port >= 0 AND port <= 65535"`
	Position int       `gorm:"type:int;uniqueIndex:idx_name_position;check:position >= 0"`
}

func (a *Automation) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	if a.Host == "" {
		return fmt.Errorf("hostname is required")
	}
	if a.Port <= 0 || a.Port > 65535 {
		return fmt.Errorf("error: Port %d is not valid", a.Port)
	}
	return nil
}
