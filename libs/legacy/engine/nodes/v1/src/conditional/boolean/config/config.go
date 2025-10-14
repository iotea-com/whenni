package booleanConditionalNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/models"

type BooleanConditionalNodeConfig struct {
	Model           models.Model              `json:"model::model" validate:"required"`
	Conditions      []models.BooleanCondition `json:"conditions" validate:"required"`
	LogicalOperator string                    `json:"logicalOperator" validate:"required,oneof=AND OR"`
}
