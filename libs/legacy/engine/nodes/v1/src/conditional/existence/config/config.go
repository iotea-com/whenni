package existenceConditionalNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/models"

type ExistenceConditionalNodeConfig struct {
	Model        models.Model `json:"model::model" validate:"required"`
	AttributeIds []string     `json:"attributeIds" validate:"required"`
}
