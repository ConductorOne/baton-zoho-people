package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// RoleResourceTypeID is the resource type ID for roles, used to check whether
// the role resource type is included in the current sync filter (see
// (*cli.ConnectorOpts).WillSyncResourceType).
const RoleResourceTypeID = "role"

// The user resource type is for all user objects from the database.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

var roleResourceType = &v2.ResourceType{
	Id:          RoleResourceTypeID,
	DisplayName: "Role",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_ROLE},
}
