package stringCompareConditionalNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/models"

type StringCompareConditionalNodeConfig struct {
	Model       models.Model              `json:"model::model" validate:"required"`
	Comparisons []models.StringComparison `json:"comparisons" validate:"required"`
}
