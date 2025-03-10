package connector

import (
	"context"
	"fmt"
	"strconv"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-zoho-people/pkg/client"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.ZohoPeopleClient
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
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
		userResource, err := parseIntoUserResource(&employeeCopy, "")
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

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, res *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant

	var userID = res.Id.Resource

	employees, _, _, err := o.client.GetEmployeeByID(ctx, userID)

	if err != nil {
		return nil, "", nil, err
	}

	for _, employee := range employees {
		if employee.Role != "" {
			departmentResource := &v2.Resource{
				Id: &v2.ResourceId{
					ResourceType: departmentResourceType.Id,
					Resource:     employee.DepartmentID,
				},
			}
			employeeCopy := employee
			userResource, _ := parseIntoUserResource(&employeeCopy, userID)
			userGrant := grant.NewGrant(departmentResource, employee.Role, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("department-grant:%s:%s:%s", departmentResource.Id.Resource, userID, employee.Role),
			}))
			grants = append(grants, userGrant)
		}
	}

	return grants, "", nil, nil
}

func parseIntoUserResource(user *client.Employee, zohoID string) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	profile := map[string]interface{}{
		"employee_id": user.EmployeeID,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"email_id":    user.EmailID,
		"zuid":        user.ZUID,
	}
	displayName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	userID := zohoID
	if user.ZohoID != 0 {
		userID = strconv.FormatInt(user.ZohoID, 10)
	}
	userTraits := []resource.UserTraitOption{
		resource.WithUserProfile(profile),
		resource.WithStatus(userStatus),
		resource.WithUserLogin(displayName),
	}

	ret, err := resource.NewUserResource(
		displayName,
		userResourceType,
		userID,
		userTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newUserBuilder(c *client.ZohoPeopleClient) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}
