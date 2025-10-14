package transformProcessingNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/models"

type TransformProcessingNodeConfig struct {
	ModelInput  models.Model      `json:"inputModel::model" validate:"required"`
	ModelOutput models.Model      `json:"outputModel::model" validate:"required"`
	Mapping     map[string]string `json:"mapping" validate:"required"`
}
