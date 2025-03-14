package connector

import (
	"context"
	"fmt"
	"strings"

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

var zohoRoles = []string{"Admin", "Team Incharge", "Team member", "Manager", "Director"}

func (o *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (o *roleBuilder) List(_ context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	for _, zohoRole := range zohoRoles {
		roleResource, err := parseIntoRoleResource(zohoRole)
		if err != nil {
			return nil, "", nil, err
		}

		resources = append(resources, roleResource)
	}

	return resources, "", nil, nil
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

func parseIntoRoleResource(zohoRole string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_name": zohoRole,
	}

	roleTraits := []resourceType.RoleTraitOption{
		resourceType.WithRoleProfile(profile),
	}

	ret, err := resourceType.NewRoleResource(
		zohoRole,
		roleResourceType,
		GetRoleID(zohoRole),
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

func GetRoleID(roleName string) string {
	return fmt.Sprintf("zoho-role_%s", strings.ToLower(strings.ReplaceAll(roleName, " ", "-")))
}
