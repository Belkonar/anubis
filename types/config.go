package types

type TargetConfig struct {
	Prefix         string                 `json:"prefix"`
	Target         string                 `json:"target"`
	RefuseFallback bool                   `json:"refuseFallback"` // If true, it'll refuse to fall back to proxy
	Endpoints      []TargetEndpointConfig `json:"endpoints"`
}

type TargetEndpointConfig struct {
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Fga    []FgaConfig `json:"fga"`
}

type AuthConfig struct {
	Issuer   string `json:"issuer"`
	Audience string `json:"audience"`
}

// FgaConfig is the configuration for an FGA relation query
//
// Uses either the object key as a param search or a specific Id
type FgaConfig struct {
	UserType   string `json:"userType"` // Defaults to "user"
	UserId     string `json:"userId"`   // Defaults to requestor subject
	Relation   string `json:"relation"`
	ObjectType string `json:"objectType"`
	ObjectKey  string `json:"objectKey"` // If this is set, it'll use a route param
	ObjectId   string `json:"objectId"`
}
