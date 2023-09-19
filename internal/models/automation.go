package models

import (
	"fmt"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"mime/multipart"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

type Automation struct {
	ID          uuid.UUID             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id,omitempty"`
	Name        string                `gorm:"type:varchar(50);unique" json:"name,omitempty"`
	URLPath     string                `gorm:"type:varchar(255);unique" json:"urlPath,omitempty"`
	Image       string                `gorm:"type:varchar(255)" json:"image,omitempty"`
	Host        string                `gorm:"type:varchar(50)" json:"host,omitempty"`
	Port        int                   `gorm:"check:port >= 0 AND port <= 65535" json:"port,omitempty"`
	Position    int                   `gorm:"type:int;unique;check:position >= 0" json:"position,omitempty,omitinput"`
	ImageFile   *multipart.FileHeader `json:"imageFile,omitempty" gorm:"-"`
	RemoveImage bool                  `json:"removeImage,omitempty" gorm:"-"`
}

func (a *Automation) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(a.Name) > 50 {
		return fmt.Errorf("name is too long, maximum length is 50 characters")
	}
	if a.URLPath == "" {
		return fmt.Errorf("urlPath is required")
	}
	if len(a.URLPath) > 255 {
		return fmt.Errorf("urlPath is too long, maximum length is 255 characters")
	}
	if len(a.Image) > 255 {
		return fmt.Errorf("image name is too long, maximum length is 255 characters")
	}
	if a.Host == "" {
		return fmt.Errorf("hostname is required")
	}
	if len(a.Host) > 50 {
		return fmt.Errorf("hostname is too long, maximum length is 50 characters")
	}
	if a.Port <= 0 || a.Port > 65535 {
		return fmt.Errorf("error: Port %d is not valid", a.Port)
	}
	if a.Position < 0 {
		return fmt.Errorf("position cannot be negative")
	}
	return nil
}
