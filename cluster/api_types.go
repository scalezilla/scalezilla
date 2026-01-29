package cluster

// APIBootstrapStatusResponse return the response to the bootstrap
// status command
type APIBootstrapStatusResponse struct {
	// Bootstrapped indicates if the cluster is boostrapped or not
	Bootstrapped bool `json:"bootstrapped"`
}

// APIBootstrapClusterRequest handle request to boostrap the cluster
type APIBootstrapClusterRequest struct {
	// Token to use to bootstrap the cluster if present
	Token string `json:"token" binding:"-"`
}

// // APIBootstrapClusterResponse handle response to boostrap the cluster
// // request
// type APIBootstrapClusterResponse struct {
// 	// AccessorID is the uuid returned with the bootstrap token
// 	AccessorID string `json:"accessor_id"`

// 	// Token is the one used to bootstrap the cluster
// 	Token string `json:"token"`

// 	// InitialToken is an boolean indicating if it's the inital bootstrap token
// 	InitialToken string `json:"initial_token"`
// }

// AclToken holds all requirements for later use in api call
type AclToken struct {
	// AccessorID is an uuid linked to then token
	AccessorID string `json:"accessor_id"`

	// Token is the one used to bootstrap the cluster
	Token string `json:"token"`

	// InitialToken is an boolean indicating if it's the inital bootstrap token
	InitialToken bool `json:"initial_token"`
}
