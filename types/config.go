package types

// FgaConfig is the configuration for an FGA relation query
//
// Uses either the object key as a param search or a specific Id
type FgaConfig struct {
	UserType   string // Defaults to "user"
	UserId     string // Defaults to requestor subject
	Relation   string
	ObjectType string
	ObjectKey  string // If this is set, it'll use a route param
	ObjectId   string
}
