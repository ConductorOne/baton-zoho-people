package connector

import (
	"context"
	"fmt"
	"strconv"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-zoho-people/pkg/client"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.ZohoPeopleClient
	// syncRoles gates the cross-type role-grant emission below: the "role"
	// resource type must not appear in emitted grants when the sync filter
	// excludes it.
	syncRoles bool
}

func (o *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, _ *v2.ResourceId, opts rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var resources []*v2.Resource

	pToken := &opts.PageToken

	bag, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, nil, err
	}
	employees, nextPageToken, _, err := o.client.ListUsers(ctx, client.PageOptions{
		PageSize:  pToken.Size,
		PageToken: pageToken,
	})

	if err != nil {
		return nil, nil, err
	}

	err = bag.Next(nextPageToken)
	if err != nil {
		return nil, nil, err
	}

	for _, employee := range employees {
		employeeCopy := employee
		userResource, err := parseIntoUserResource(&employeeCopy, "")
		if err != nil {
			return nil, nil, err
		}

		resources = append(resources, userResource)
	}

	next, err := bag.Marshal()
	if err != nil {
		return nil, nil, err
	}

	return resources, &rs.SyncOpResults{NextPageToken: next}, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants emits role-assignment grants for the user, as a cross-type
// optimization: the Zoho "employee" API response already includes the user's
// role, so the user builder emits grants against the role resource type
// instead of requiring the role builder to make an additional call per user.
//
// This must not run when the "role" resource type is excluded from the
// current sync filter - the connector should never emit grants referencing a
// resource type it isn't syncing. Users have no entitlements of their own
// (see Entitlements above), so skipping entirely in that case is correct.
func (o *userBuilder) Grants(ctx context.Context, res *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	if !o.syncRoles {
		return nil, nil, nil
	}

	var grants []*v2.Grant

	var userID = res.Id.Resource

	employees, _, _, err := o.client.GetEmployeeByID(ctx, userID)

	if err != nil {
		return nil, nil, err
	}

	for _, employee := range employees {
		roleName := employee.Role
		if roleName != "" {
			roleResource := &v2.Resource{
				Id: &v2.ResourceId{
					ResourceType: roleResourceType.Id,
					Resource:     GetRoleID(roleName),
				},
			}
			employeeCopy := employee
			userResource, _ := parseIntoUserResource(&employeeCopy, userID)
			userGrant := grant.NewGrant(roleResource, "assigned", userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("role-grant:%s:%s:%s", GetRoleID(roleName), userID, "assigned"),
			}))
			grants = append(grants, userGrant)
		}
	}

	return grants, nil, nil
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
	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(displayName),
	}

	ret, err := rs.NewUserResource(
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

func newUserBuilder(c *client.ZohoPeopleClient, syncRoles bool) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
		syncRoles:    syncRoles,
	}
}
