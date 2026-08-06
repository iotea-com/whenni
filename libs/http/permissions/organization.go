package gruentpermissions

var NamespaceOrganizations Namespace = "organizations"
var NamespaceOrganizationApiKeys Namespace = "organization-api-keys"
var NamespaceSpaces Namespace = "spaces"
var NamespaceMembers Namespace = "members"
var NamespacePermissions Namespace = "permissions"

// Default permission set for organization members
var DefaultMemberOrganizationPermissions = []string{
	// Organization
	Marshal(NamespaceOrganizations, ActionGet),
	// Marshal(NamespaceOrganizations, ActionUpdate),
	// Marshal(NamespaceOrganizations, ActionDelete),

	// API keys
	Marshal(NamespaceOrganizationApiKeys, ActionList),
	// Marshal(NamespaceOrganizationApiKeys, ActionCreate),
	// Marshal(NamespaceOrganizationApiKeys, ActionDelete),

	// Spaces
	Marshal(NamespaceSpaces, ActionList),
	Marshal(NamespaceSpaces, ActionCreate),
	Marshal(NamespaceSpaces, ActionGet),
	Marshal(NamespaceSpaces, ActionUpdate),
	Marshal(NamespaceSpaces, ActionDelete),

	// Members
	Marshal(NamespaceMembers, ActionList),
	// Marshal(NamespaceMembers, ActionCreate),
	Marshal(NamespaceMembers, ActionGet),
	// Marshal(NamespaceMembers, ActionUpdate),
	// Marshal(NamespaceMembers, ActionDelete),

	// Permission sets
	Marshal(NamespacePermissions, ActionList),
	// Marshal(NamespacePermissions, ActionCreate),
	Marshal(NamespacePermissions, ActionGet),
	// Marshal(NamespacePermissions, ActionUpdate),
	// Marshal(NamespacePermissions, ActionDelete),
}
