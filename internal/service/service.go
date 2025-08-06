package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/yashzod/splitlinks/internal/model"
	"gorm.io/gorm"

	"gorm.io/datatypes"
)

func pickVariant(variants []model.Variant) model.Variant {
	total := 0
	for _, v := range variants {
		total += v.Weight
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(total)

	sum := 0
	for _, v := range variants {
		sum += v.Weight
		if r < sum {
			return v
		}
	}
	return variants[0]
}

func GetRedirectLink(db *gorm.DB, slug string) (string, error) {
	var exp model.Experiment
	err := db.Where("slug = ?", slug).First(&exp).Error
	if err != nil {
		return "", err
	}
	experimentID := exp.ID

	var variants []model.Variant
	err = db.Where("experiment_id = ?", experimentID).Find(&variants).Error
	selected_variant := pickVariant(variants)
	link := selected_variant.URL
	if link == "" {
		return "", errors.New("slug not found")
	}
	return link, err
}

func GetExperiment(query map[string]string, db *gorm.DB) ([]interface{}, error) {
	var exps []interface{}
	slug := query["slug"]
	if slug != "" {
		err := db.Where("slug = ?", slug).Find(&exps).Error
		return exps, err
	}
	user_id := query["user_id"]
	err := db.Where("CreatedID = ?", user_id).Find(&exps).Error
	return exps, err

}

func CreateExperiment(db *gorm.DB, data map[string]interface{}) (string, error) {
	expData, ok := data["experiment"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid or missing experiment data")
	}

	variantsData, ok := data["variants"].([]interface{})
	if !ok || len(variantsData) == 0 {
		return "", errors.New("invalid or missing variants data")
	}

	experimentID := uuid.New()
	slug, err := generateUniqueSlug(db)
	if err != nil {
		return "", err
	}
	println("slug exp created")

	fmt.Printf("slug: %s\n", slug)

	experiment := model.Experiment{
		ID:        experimentID,
		Slug:      slug,
		Name:      expData["name"].(string),
		CreatedAt: time.Now(),
	}

	if md, ok := expData["metadata"]; ok {
		jsonVal, err := json.Marshal(md)
		if err == nil {
			experiment.Metadata = datatypes.JSON(jsonVal)
		}
	}

	if err := db.Create(&experiment).Error; err != nil {
		return "", err
	}

	for _, v := range variantsData {
		vMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		var targeting datatypes.JSON

		if rawTargeting, ok := vMap["targeting"]; ok {
			jsonVal, err := json.Marshal(rawTargeting)
			if err != nil {
				fmt.Println("Error marshalling targeting:", err)
			} else {
				targeting = datatypes.JSON(jsonVal)
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
			return "", err
		}
	}

	return slug, nil
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
