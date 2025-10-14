package apitest

import (
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/util"
)

func SetupExpectOrgApiKey(dbMocks handlerDbMocks) {
	orgId, _ := id.Generator.NewOrganizationId()
	orgPermSetId, _ := id.Generator.NewPermissionSetId()
	now := util.GetCurrentTime()

	dbMocks.Server.APIKey.Expect(
		dbMocks.Client.APIKey.FindUnique(
			db.APIKey.ID.Equals("tea_fake"),
		).With(
			db.APIKey.Organization.Fetch().With(
				db.Organization.Spaces.Fetch(),
			),
			db.APIKey.OrganizationPermissionSet.Fetch(),
			db.APIKey.SpacePermissionSet.Fetch(),
		),
	).Returns(db.APIKeyModel{
		InnerAPIKey: db.InnerAPIKey{
			ID:                          "tea_fake",
			OrganizationID:              *orgId,
			OrganizationPermissionSetID: orgPermSetId,
			CreatedAt:                   now,
			CreatedBy:                   "test",
		},
		RelationsAPIKey: db.RelationsAPIKey{
			Organization: &db.OrganizationModel{
				InnerOrganization: db.InnerOrganization{
					ID:        *orgId,
					Name:      "Test Org",
					CreatedAt: now,
					UpdatedAt: now,
					CreatedBy: "test",
					UpdatedBy: "test",
				},
				RelationsOrganization: db.RelationsOrganization{
					Spaces: []db.SpaceModel{},
				},
			},
			OrganizationPermissionSet: &db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					ID:             *orgPermSetId,
					OrganizationID: *orgId,
					Name:           "test",
					Permissions:    []string{"*"},
					CreatedAt:      now,
					UpdatedAt:      now,
					CreatedBy:      "test",
					UpdatedBy:      "test",
				},
			},
		},
	})
}

func SetupExpectSpaceApiKey(dbMocks handlerDbMocks) {
	orgId, _ := id.Generator.NewOrganizationId()
	spaceId, _ := id.Generator.NewSpaceId()
	// orgPermSetId, _ := id.Generator.NewPermissionSetId()
	spacePermSetId, _ := id.Generator.NewPermissionSetId()
	now := util.GetCurrentTime()

	dbMocks.Server.APIKey.Expect(
		dbMocks.Client.APIKey.FindUnique(
			db.APIKey.ID.Equals("tea_fake"),
		).With(
			db.APIKey.Organization.Fetch().With(
				db.Organization.Spaces.Fetch(),
			),
			db.APIKey.OrganizationPermissionSet.Fetch(),
			db.APIKey.SpacePermissionSet.Fetch(),
		),
	).Returns(db.APIKeyModel{
		InnerAPIKey: db.InnerAPIKey{
			ID:                          "tea_fake",
			OrganizationID:              *orgId,
			SpaceID:                     spaceId,
			OrganizationPermissionSetID: nil,
			SpacePermissionSetID:        spacePermSetId,
			CreatedAt:                   now,
			CreatedBy:                   "test",
		},
		RelationsAPIKey: db.RelationsAPIKey{
			Organization: &db.OrganizationModel{
				InnerOrganization: db.InnerOrganization{
					ID: *orgId,
				},
				RelationsOrganization: db.RelationsOrganization{
					Spaces: []db.SpaceModel{
						{
							InnerSpace: db.InnerSpace{
								ID: *spaceId,
							},
						},
					},
				},
			},
			SpacePermissionSet: &db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					ID:             *spacePermSetId,
					Permissions:    []string{"*"},
					OrganizationID: *orgId,
					SpaceID:        spaceId,
					Name:           "test",
					CreatedAt:      now,
					UpdatedAt:      now,
					CreatedBy:      "test",
					UpdatedBy:      "test",
				},
			},
		},
	})
}
