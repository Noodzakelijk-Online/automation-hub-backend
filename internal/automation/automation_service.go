package automation

import (
	"automation-hub-backend/internal/config"
	"automation-hub-backend/internal/model"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	FindByID(id uuid.UUID) (*model.Automation, error)
	Create(automation *model.Automation) (*model.Automation, error)
	Update(automation *model.Automation) (*model.Automation, error)
	Delete(id uuid.UUID) error
	FindAll() ([]*model.Automation, error)
	SwapOrder(id1 uuid.UUID, id2 uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func DefaultService() Service {
	repo := DefaultRepository()
	return NewService(repo)
}

func (s *service) FindByID(id uuid.UUID) (*model.Automation, error) {
	return s.repo.FindByID(id)
}

func (s *service) Create(automation *model.Automation) (*model.Automation, error) {
	automation.ID = uuid.UUID{} // reset ID

	if automation.ImageFile != nil {
		newFileName, err := s.processImageFile(automation.ImageFile)
		if err != nil {
			return nil, err
		}
		automation.Image = newFileName
	}

	if err := automation.Validate(); err != nil {
		return nil, err
	}

	maxPosition, err := s.repo.MaxPosition()
	if err != nil {
		return nil, err
	}
	automation.Position = maxPosition + 1

	return s.repo.Create(automation)
}

func (s *service) Update(automation *model.Automation) (*model.Automation, error) {
	currentAutomation, err := s.repo.FindByID(automation.ID)
	if err != nil {
		return nil, err
	}

	automation.Position = currentAutomation.Position

	if automation.ImageFile != nil {
		newFileName, err := s.processImageFile(automation.ImageFile)
		if err != nil {
			return nil, err
		}
		if err := s.deleteImage(currentAutomation.Image); err != nil {
			return nil, err
		}
		automation.Image = newFileName
	} else if automation.RemoveImage {
		if err := s.deleteImage(currentAutomation.Image); err != nil {
			return nil, err
		}
		automation.Image = ""
	} else {
		automation.Image = currentAutomation.Image
	}

	if err := automation.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Update(automation)
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) FindAll() ([]*model.Automation, error) {
	return s.repo.FindAll()
}

func (s *service) SwapOrder(id1 uuid.UUID, id2 uuid.UUID) error {
	return s.repo.Transaction(func(tx *gorm.DB) error {
		automation1, err := s.repo.FindByID(id1)
		if err != nil {
			return err
		}
		automation2, err := s.repo.FindByID(id2)
		if err != nil {
			return err
		}

		pos1 := automation1.Position
		pos2 := automation2.Position

		maxPosition, err := s.repo.MaxPosition()
		if err != nil {
			return err
		}
		tempPosition := maxPosition + 1

		automation1.Position = tempPosition
		if err := tx.Save(automation1).Error; err != nil {
			return err
		}

		automation2.Position = pos1
		if err := tx.Save(automation2).Error; err != nil {
			return err
		}

		automation1.Position = pos2
		if err := tx.Save(automation1).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *service) processImageFile(file *multipart.FileHeader) (string, error) {
	if file.Size > config.AppConfig.ImageMaxSize {
		return "", fmt.Errorf("image is too large (%d). Max size is %d Mb", file.Size, config.AppConfig.ImageMaxSize)
	}

	ext := filepath.Ext(file.Filename)
	fmt.Printf("Filename: %s, Extracted Extension: %s\n", file.Filename, ext)

	if !contains(config.AppConfig.ImageExtensions, ext) {
		return "", fmt.Errorf("invalid image extension. Allowed extensions are: %v", config.AppConfig.ImageExtensions)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}(src)

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)
	if !strings.HasPrefix(fileType, "image/") {
		return "", fmt.Errorf("file is not an image")
	}
	mimeSuffix := strings.TrimPrefix(fileType, "image/")
	if !contains(config.AppConfig.ImageExtensions, "."+mimeSuffix) {
		return "", fmt.Errorf("mismatch between file extension and MIME type")
	}

	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	_, _, err = image.Decode(src)
	if err != nil {
		//return "", fmt.Errorf("corrupted image: %v", err)
	}

	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	newFileName := uuid.New().String() + ext
	dst, err := os.Create(config.AppConfig.ImageSaveDir + "/" + newFileName)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			fmt.Printf("Failed to close file %s: %v", dst.Name(), err)
		}
	}(dst)

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func (s *service) deleteImage(imageName string) error {
	if imageName == "" {
		return nil
	}
	imagePath := config.AppConfig.ImageSaveDir + "/" + imageName
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(imagePath)
}

func contains(slice []string, str string) bool {
	str = strings.ToLower(str)
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
