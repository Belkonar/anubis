package types

type TargetConfig struct {
    Id             string                 `json:"id"`
    Prefix         string                 `json:"prefix"`
    Target         string                 `json:"target"`
    RefuseFallback bool                   `json:"refuseFallback"` // If true, it'll refuse to fall back to proxy
    FgaModel       FgaModelConfig         `json:"fgaModel"`
    Endpoints      []TargetEndpointConfig `json:"endpoints"`
}

type TargetEndpointConfig struct {
    Id     string      `json:"id"`
    Path   string      `json:"path"`
    Method string      `json:"method"`
    Fga    []FgaConfig `json:"fga"`
}

type AuthConfig struct {
    Id       string `json:"id"`
    Issuer   string `json:"issuer"`
    Audience string `json:"audience"`
}

// FgaClusterConfig is the configuration for an FGA cluster
// Represented in the database as /fga-cluster/{id}
// May put auth in here if it can be across models.
type FgaClusterConfig struct {
    Id       string `json:"id"`
    Endpoint string `json:"endpoint"`
}

type FgaModelConfig struct {
    Id           string `json:"id"`
    FgaClusterId string `json:"fgaClusterId"`
}

// FgaConfig is the configuration for an FGA relation query
//
// Uses either the object key as a param search or a specific Id
type FgaConfig struct {
    Id         string `json:"id"`
    UserType   string `json:"userType"` // Defaults to "user"
    UserId     string `json:"userId"`   // Defaults to requester subject
    Relation   string `json:"relation"`
    ObjectType string `json:"objectType"`
    ObjectKey  string `json:"objectKey"` // If this is set, it'll use a route param
    ObjectId   string `json:"objectId"`
}
