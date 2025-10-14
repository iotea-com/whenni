package apitest

import "github.com/iotea-com/iotea/prisma/db"

func SetupExpectJwt(dbMocks handlerDbMocks, organizationId string, userId string) {
	dbMocks.Server.OrganizationMember.Expect(
		dbMocks.Client.OrganizationMember.FindUnique(
			db.OrganizationMember.OrganizationIDUserID(
				db.OrganizationMember.OrganizationID.Equals(organizationId),
				db.OrganizationMember.UserID.Equals(userId),
			),
		).With(
			db.OrganizationMember.User.Fetch(),
			db.OrganizationMember.OrganizationPermissionSet.Fetch(),
		),
	).Returns(db.OrganizationMemberModel{
		InnerOrganizationMember: db.InnerOrganizationMember{
			OrganizationID: organizationId,
			UserID:         userId,
			Role:           db.OrganizationRoleAdmin,
		},
		RelationsOrganizationMember: db.RelationsOrganizationMember{
			OrganizationPermissionSet: &db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					Permissions: []string{"*"},
				},
			},
		},
	})
}
