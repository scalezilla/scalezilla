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

// AclToken holds all requirements for later use in api call
type AclToken struct {
	// AccessorID is an uuid linked to then token
	AccessorID string `json:"accessor_id"`

	// Token is the one used to bootstrap the cluster
	Token string `json:"token"`

	// InitialToken is an boolean indicating if it's the inital bootstrap token
	InitialToken bool `json:"initial_token"`
}

// APINodesListRequest handle request to list cluster nodes
type APINodesListRequest struct {
	// Kind is cluster node kind to return when present
	Kind string `form:"kind" binding:"-"`
}

// APINodesListResponse handle response to list cluster nodes
type APINodesListResponse struct {
	// ID is the node id
	ID string `json:"id"`

	// Name is the node name
	Name string `json:"name"`

	// Address is the node host address
	Address string `json:"address"`

	// Kind is the node kind
	Kind string `json:"kind"`

	// Leader is set to true when the node is the leader
	Leader bool `json:"leader"`

	// NodePool is the node pool
	NodePool string `json:"nodePool"`
}
