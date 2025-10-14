package tagsRemove

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	SpaceId     string `validate:"required"`
	TagId       string `validate:"required"`
	SubjectId   string `validate:"required"`
}

type Output struct{}

// @Summary Remove a tag
// @Description Given a tag ID and subject ID, removes the tag from the subject.
// @Tags tags
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   tagId	query	string	true	"Tag ID"
// @Param   subjectId	query	string	true	"Subject ID"
// @Success 200
// @Router /tags/{tagId}/remove [patch]
// @Security ApiKeyAuth
func Handler(ctx *fiber.Ctx) error {
	request, err := parse(ctx)
	if err != nil {
		return err
	}

	err = validate(request)
	if err != nil {
		return err
	}

	err = contextValidate(request)
	if err != nil {
		return err
	}

	output, err := execute(request)
	if err != nil {
		return err
	}

	respond(request, output)

	go action(output)

	return nil
}
