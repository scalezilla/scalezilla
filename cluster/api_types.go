package cluster

import "time"

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
	// Kind is cluster node kind to return when present for GET requests
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

// APIDeploymentApplyRequest handle request to create deployment
type APIDeploymentApplyRequest struct {
	// HCLContent is the payload to use to create a new deployment
	HCLContent string `json:"hcl_content" binding:"required"`
}

// APIGenericResponse handle generic success response
type APIGenericResponse struct {
	Message string `json:"message"`
}

// APIPodsListRequest handle request to list pods
type APIPodsListRequest struct {
	// Namespace is the namespace to use to list pods
	Namespace string `form:"namespace,default=default"`
}

// APIPodsListResponse handle response list pods
type APIPodsListResponse struct {
	// Namespace is the namespace to use to list pods
	Namespace string `json:"namespace"`

	// ID is the container id
	ID string `json:"id"`

	// Image is the container image
	Image string `json:"image"`

	// PID is the container pid
	PID uint32 `json:"pid"`

	// Runtime is the container runtime
	Runtime string `json:"runtime"`

	// Status is the container status
	Status string `json:"status"`

	// CreatedAt is the container creation date
	CreatedAt time.Time `json:"created_at"`
}
