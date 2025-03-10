package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceType "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-zoho-people/pkg/client"
)

type departmentBuilder struct {
	resourceType *v2.ResourceType
	client       *client.ZohoPeopleClient
}

var roles = []string{"Admin", "Team member", "Manager", "Director", "Team Incharge", "Team Lead"}

func (o *departmentBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return departmentResourceType
}

func (o *departmentBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	bag, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, "", nil, err
	}
	departments, nextPageToken, _, err := o.client.ListDepartments(ctx, client.PageOptions{
		PageSize:  pToken.Size,
		PageToken: pageToken,
	})

	if err != nil {
		return nil, "", nil, err
	}

	err = bag.Next(nextPageToken)
	if err != nil {
		return nil, "", nil, err
	}

	for _, department := range departments {
		departmentCopy := department
		userResource, err := parseIntoDepartmentResource(ctx, &departmentCopy, nil)
		if err != nil {
			return nil, "", nil, err
		}

		resources = append(resources, userResource)
	}

	nextPageToken, err = bag.Marshal()
	if err != nil {
		return nil, "", nil, err
	}

	return resources, nextPageToken, nil, nil
}

func (o *departmentBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement
	for _, role := range roles {
		assigmentOptions := []entitlement.EntitlementOption{
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Zoho role %s for department %s", role, resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s Department %s", resource.DisplayName, role)),
		}

		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, role, assigmentOptions...))
	}

	return entitlements, "", nil, nil
}

func (o *departmentBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant

	var departmentID = resource.Id.Resource

	departments, _, _, err := o.client.GetDepartmentByID(ctx, departmentID)

	if err != nil {
		return nil, "", nil, err
	}

	for _, department := range departments {
		if department.DepartmentLeadID != "" {
			userResource, _ := parseIntoUserResourceFromDepartment(&department)
			leadGrant := grant.NewGrant(resource, "Team Lead", userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("department-grant:%s:%s:%s", resource.Id.Resource, department.DepartmentLeadID, "Team Lead"),
			}))
			grants = append(grants, leadGrant)
		}
	}

	return grants, "", nil, nil
}

func parseIntoUserResourceFromDepartment(department *client.Department) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	profile := map[string]interface{}{
		"zoho_id": department.DepartmentLeadID,
		"email":   department.DepartmentLeadMail,
	}
	displayName := department.DepartmentLead
	userTraits := []resourceType.UserTraitOption{
		resourceType.WithUserProfile(profile),
		resourceType.WithStatus(userStatus),
		resourceType.WithUserLogin(displayName),
	}

	userResource, err := resourceType.NewUserResource(
		department.DepartmentLead,
		userResourceType,
		department.DepartmentLeadID,
		userTraits,
		resourceType.WithParentResourceID(nil),
	)

	if err != nil {
		return nil, err
	}

	return userResource, nil
}

func parseIntoDepartmentResource(_ context.Context, department *client.Department, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"department_id": department.ZohoID,
		"department":    department.Department,
		"email":         department.MailAlias,
	}

	groupTraits := []resourceType.GroupTraitOption{
		resourceType.WithGroupProfile(profile),
	}

	displayName := department.Department

	if department.ParentDepartmentID != "" {
		parentResource, _ := resourceType.NewResourceID(departmentResourceType, department.ParentDepartmentID)
		parentResourceID = parentResource
	}

	ret, err := resourceType.NewGroupResource(
		displayName,
		departmentResourceType,
		department.ZohoID,
		groupTraits,
		resourceType.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil

}

func newDepartmentBuilder(c *client.ZohoPeopleClient) *departmentBuilder {
	return &departmentBuilder{
		resourceType: departmentResourceType,
		client:       c,
	}
}
