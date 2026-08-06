package gruentpermissions

var NamespaceSpaceApiKeys Namespace = "space-api-keys"
var NamespaceChannels Namespace = "channels"
var NamespaceChannelExecutions Namespace = "channel-executions"
var NamespaceAnalytics Namespace = "analytics"
var NamespaceThings Namespace = "things"
var NamespaceModels Namespace = "models"
var NamespaceData Namespace = "data"
var NamespacePolicies Namespace = "policies"
var NamespaceCertificates Namespace = "certificates"
var NamespaceEnvironments Namespace = "environments"
var NamespaceSecrets Namespace = "secrets"

// Default permission set for space members
var DefaultMemberSpacePermissions = []string{
	// API Keys
	Marshal(NamespaceSpaceApiKeys, ActionList),
	// Marshal(NamespaceSpaceApiKeys, ActionCreate),
	// Marshal(NamespaceSpaceApiKeys, ActionDelete),

	// Channels
	Marshal(NamespaceChannels, ActionList),
	Marshal(NamespaceChannels, ActionCreate),
	Marshal(NamespaceChannels, ActionGet),
	Marshal(NamespaceChannels, ActionUpdate),
	// Marshal(NamespaceChannels, ActionDelete),

	// Channel executions
	Marshal(NamespaceChannelExecutions, ActionList),
	Marshal(NamespaceChannelExecutions, ActionGet),

	// // Analytics
	// Marshal(NamespaceAnalytics, ActionList),
	// Marshal(NamespaceAnalytics, ActionCreate),
	// Marshal(NamespaceAnalytics, ActionGet),
	// Marshal(NamespaceAnalytics, ActionUpdate),
	// Marshal(NamespaceAnalytics, ActionDelete),

	// Things
	Marshal(NamespaceThings, ActionList),
	Marshal(NamespaceThings, ActionCreate),
	Marshal(NamespaceThings, ActionGet),
	Marshal(NamespaceThings, ActionUpdate),
	// Marshal(NamespaceThings, ActionDelete),

	// Models
	Marshal(NamespaceModels, ActionList),
	Marshal(NamespaceModels, ActionCreate),
	Marshal(NamespaceModels, ActionGet),
	Marshal(NamespaceModels, ActionUpdate),
	// Marshal(NamespaceModels, ActionDelete),

	// // Data
	// Marshal(NamespaceData, ActionList),
	// Marshal(NamespaceData, ActionCreate),
	// Marshal(NamespaceData, ActionGet),
	// Marshal(NamespaceData, ActionUpdate),
	// Marshal(NamespaceData, ActionDelete),

	// Policies
	Marshal(NamespacePolicies, ActionList),
	// Marshal(NamespacePolicies, ActionCreate),
	// Marshal(NamespacePolicies, ActionGet),
	// Marshal(NamespacePolicies, ActionUpdate),
	// Marshal(NamespacePolicies, ActionDelete),

	// Certificates
	Marshal(NamespaceCertificates, ActionList),
	// Marshal(NamespaceCertificates, ActionCreate),
	// Marshal(NamespaceCertificates, ActionGet),
	// Marshal(NamespaceCertificates, ActionUpdate),
	// Marshal(NamespaceCertificates, ActionDelete),

	// Environments
	Marshal(NamespaceEnvironments, ActionList),
	// Marshal(NamespaceEnvironments, ActionCreate),
	// Marshal(NamespaceEnvironments, ActionGet),
	// Marshal(NamespaceEnvironments, ActionUpdate),
	// Marshal(NamespaceEnvironments, ActionDelete),

	// Secrets
	Marshal(NamespaceSecrets, ActionList),
	// Marshal(NamespaceSecrets, ActionCreate),
	// Marshal(NamespaceSecrets, ActionUpdate),
	// Marshal(NamespaceSecrets, ActionDelete),
}
