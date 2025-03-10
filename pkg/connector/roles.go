package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	resourceType "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-zoho-people/pkg/client"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.ZohoPeopleClient
}

func (o *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (o *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	bag, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, "", nil, err
	}
	employees, nextPageToken, _, err := o.client.ListUsers(ctx, client.PageOptions{
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

	for _, employee := range employees {
		employeeCopy := employee
		roleResource, err := parseIntoRoleResource(&employeeCopy)
		if err != nil {
			return nil, "", nil, err
		}

		resources = append(resources, roleResource)
	}

	nextPageToken, err = bag.Marshal()
	if err != nil {
		return nil, "", nil, err
	}

	return resources, nextPageToken, nil, nil
}

func (o *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Zoho role %s", resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("assigned role %s", resource.DisplayName)),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, "assigned", assigmentOptions...))

	return entitlements, "", nil, nil
}

func (o *roleBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoRoleResource(employee *client.Employee) (*v2.Resource, error) {
	roleID := employee.RoleID
	roleName := employee.Role

	profile := map[string]interface{}{
		"role_id":   roleID,
		"role_name": roleName,
	}

	roleTraits := []resourceType.RoleTraitOption{
		resourceType.WithRoleProfile(profile),
	}

	ret, err := resourceType.NewRoleResource(
		roleName,
		roleResourceType,
		roleID,
		roleTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newRoleBuilder(c *client.ZohoPeopleClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}
