package organizationsMembersList

import (
	"fmt"
	"strings"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.OrgId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	// Count members in organization (based on filter)
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "List members in organization")
	organizationMembers, err := sqlc.Queries.ListOrganizationMembers(dbCtx, request.Input.OrgId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing members in organization with ID %s in the database: %s", request.Input.OrgId, err)),
		)
		dbSpan.End()

		return nil, err
	}
	filteredMembers := make([]sqldb.ListOrganizationMembersRow, 0, len(organizationMembers))
	filter := strings.TrimSpace(strings.ToLower(request.Input.Filter))
	for _, member := range organizationMembers {
		userEmail := ""
		if member.UserEmail != nil {
			userEmail = *member.UserEmail
		}
		if filter == "" || strings.Contains(strings.ToLower(userEmail), filter) {
			filteredMembers = append(filteredMembers, member)
		}
	}
	membersCount := len(filteredMembers)

	start := (request.Input.Page - 1) * request.Input.ResultsPerPage
	if start > membersCount {
		start = membersCount
	}
	end := start + request.Input.ResultsPerPage
	if end > membersCount {
		end = membersCount
	}

	dbSpan.SetAttributes(attribute.Int("members_count", membersCount))
	dbSpan.End()

	output := Output{
		OrganizationMembers: filteredMembers[start:end],
		Page:                request.Input.Page,
		TotalPages:          (membersCount + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:        membersCount,
		ResultsPerPage:      request.Input.ResultsPerPage,
	}

	return &output, nil
}
