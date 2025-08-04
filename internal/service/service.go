package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/yashzod/splitlinks/internal/model"
	"gorm.io/gorm"
)

var slug_map = map[string]string{
	"abc": "https://google.com",
}

func GetRedirectLink(slug string) (string, error) {
	link := slug_map[slug]
	if link == "" {
		return "", errors.New("slug not found")
	}
	return link, nil
}

func CreateExperiment(data map[string]interface{}, db *gorm.DB) error {
	expData, ok := data["experiment"].(map[string]interface{})
	if !ok {
		return errors.New("invalid or missing experiment data")
	}

	variantsData, ok := data["variants"].([]interface{})
	if !ok || len(variantsData) == 0 {
		return errors.New("invalid or missing variants data")
	}

	experimentID := uuid.New()
	slug, err := generateUniqueSlug(db)
	if err != nil {
		return err
	}

	experiment := model.Experiment{
		ID:        experimentID,
		Slug:      slug,
		Name:      expData["name"].(string),
		CreatedAt: time.Now(),
	}

	if md, ok := expData["metadata"].(map[string]interface{}); ok {
		metadata := make(map[string]string)
		for k, v := range md {
			metadata[k] = fmt.Sprintf("%v", v)
		}
		experiment.Metadata = metadata
	}

	if err := db.Create(&experiment).Error; err != nil {
		return err
	}

	for _, v := range variantsData {
		vMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		targeting := make(map[string][]string)
		if rawTargeting, ok := vMap["targeting"].(map[string]interface{}); ok {
			for k, val := range rawTargeting {
				if arr, ok := val.([]interface{}); ok {
					strs := make([]string, len(arr))
					for i, item := range arr {
						strs[i] = fmt.Sprintf("%v", item)
					}
					targeting[k] = strs
				}
			}
		}

		variant := model.Variant{
			ID:           uuid.New(),
			ExperimentID: experiment.ID,
			Name:         vMap["name"].(string),
			URL:          vMap["url"].(string),
			Weight:       int(vMap["weight"].(float64)), // JSON numbers are float64
			Targeting:    targeting,
		}

		if err := db.Create(&variant).Error; err != nil {
			return err
		}
	}

	return nil
}

func generateUniqueSlug(db *gorm.DB) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 10; i++ {
		slug := randomString(6, charset)
		var count int64
		db.Model(&model.Experiment{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			return slug, nil
		}
	}
	return "", errors.New("could not generate unique slug")
}

func randomString(length int, charset string) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}
