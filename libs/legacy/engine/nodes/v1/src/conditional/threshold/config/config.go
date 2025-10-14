package thresholdConditionalNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/models"

type ThresholdConditionalNodeConfig struct {
	Model           models.Model       `json:"model::model" validate:"required"`
	Conditions      []models.Condition `json:"conditions" validate:"required"`
	LogicalOperator string             `json:"logicalOperator" validate:"required,oneof=AND OR"`
}
