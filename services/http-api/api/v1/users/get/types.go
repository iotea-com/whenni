package get

import (
	"time"

	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type User struct {
	ID            string               `json:"id"`
	Name          *string              `json:"name,omitempty"`
	Email         *string              `json:"email,omitempty"`
	EmailVerified *time.Time           `json:"emailVerified,omitempty"`
	Image         *string              `json:"image,omitempty"`
	Organizations []OrganizationMember `json:"organizations,omitempty"`
}

type OrganizationMember struct {
	ID                          string        `json:"id"`
	OrganizationID              string        `json:"organizationId"`
	UserID                      string        `json:"userId"`
	Role                        string        `json:"role"`
	OrganizationPermissionSetID string        `json:"organizationPermissionSetId"`
	CreatedAt                   time.Time     `json:"createdAt"`
	UpdatedAt                   time.Time     `json:"updatedAt"`
	Organization                *Organization `json:"organization,omitempty"`
}

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
	Spaces    []Space   `json:"spaces,omitempty"`
}

type Space struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
	CreatedBy      string    `json:"createdBy"`
	UpdatedAt      time.Time `json:"updatedAt"`
	UpdatedBy      string    `json:"updatedBy"`
}

func toUser(user sqldb.AppUser, members []sqldb.ListUserOrganizationMembersRow, spacesByOrg map[string][]sqldb.AppSpace) *User {
	var emailVerified *time.Time
	if user.EmailVerified.Valid {
		verifiedAt := user.EmailVerified.Time
		emailVerified = &verifiedAt
	}

	organizations := make([]OrganizationMember, 0, len(members))
	for _, member := range members {
		spaces := make([]Space, 0, len(spacesByOrg[member.OrganizationID]))
		for _, space := range spacesByOrg[member.OrganizationID] {
			spaces = append(spaces, Space{
				ID:             space.ID,
				OrganizationID: space.OrganizationID,
				Name:           space.Name,
				CreatedAt:      space.CreatedAt,
				CreatedBy:      space.CreatedBy,
				UpdatedAt:      space.UpdatedAt,
				UpdatedBy:      space.UpdatedBy,
			})
		}

		organizations = append(organizations, OrganizationMember{
			ID:                          member.ID,
			OrganizationID:              member.OrganizationID,
			UserID:                      member.UserID,
			Role:                        string(member.Role),
			OrganizationPermissionSetID: member.OrganizationPermissionSetID,
			CreatedAt:                   member.CreatedAt,
			UpdatedAt:                   member.UpdatedAt,
			Organization: &Organization{
				ID:        member.OrganizationID,
				Name:      member.OrganizationName,
				CreatedAt: member.OrganizationCreatedAt,
				CreatedBy: member.OrganizationCreatedBy,
				UpdatedAt: member.OrganizationUpdatedAt,
				UpdatedBy: member.OrganizationUpdatedBy,
				Spaces:    spaces,
			},
		})
	}

	return &User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: emailVerified,
		Image:         user.Image,
		Organizations: organizations,
	}
}
