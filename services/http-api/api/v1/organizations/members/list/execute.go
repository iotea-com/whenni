package organizationsMembersList

import (
	"fmt"
	"strconv"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.OrgId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	// Count members in organization (based on filter)
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count members in space")
	var membersCountResponse []struct {
		Count db.RawString `json:"members_count"`
	}
	countQuery := `
		SELECT count(*) as members_count 
		FROM app.organization_members 
		JOIN app.users ON users.id = organization_members."userId"
		WHERE organization_members."organizationId" = $1`

	queryParams := []interface{}{request.Input.OrgId}
	paramCount := 1

	if request.Input.Filter != "" {
		paramCount++
		countQuery += fmt.Sprintf(` AND (
			LOWER(users.email) LIKE LOWER($%d)
		)`, paramCount)
		queryParams = append(queryParams, "%"+request.Input.Filter+"%")
	}

	err := prisma.Client.Prisma.QueryRaw(countQuery, queryParams...).Exec(dbCtx, &membersCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting members in organization in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	membersCount, err := strconv.ParseInt(string(membersCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", membersCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("members_count", membersCount),
	)
	dbSpan.End()

	// List members in organization (based on filter)
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List members in organization")
	var conditions []db.OrganizationMemberWhereParam = []db.OrganizationMemberWhereParam{
		db.OrganizationMember.OrganizationID.Equals(request.Input.OrgId),
	}

	if request.Input.Filter != "" {
		conditions = append(conditions, db.OrganizationMember.Or(
			db.OrganizationMember.User.Where(
				db.User.Email.Mode(db.QueryModeInsensitive),
				db.User.Email.Contains(request.Input.Filter),
			),
		))
	}

	organizationMembers, err := prisma.Client.OrganizationMember.FindMany(
		conditions...,
	).With(
		db.OrganizationMember.User.Fetch(), // .Omit(db.User.Password.Field()) is causing problems here - we will still filter out passwords before returning
	).Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil && err.Error() != "ErrNotFound" {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing members in organization with ID %s in the database: %s", request.Input.OrgId, err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Filter out passwords manually
	for _, member := range organizationMembers {
		user := member.User()
		user.InnerUser.Password = nil
	}

	output := Output{
		OrganizationMembers: organizationMembers,
		Page:                request.Input.Page,
		TotalPages:          (int(membersCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:        int(membersCount),
		ResultsPerPage:      request.Input.ResultsPerPage,
	}

	return &output, nil
}
