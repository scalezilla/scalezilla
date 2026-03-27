package cluster

import (
	"regexp"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/scalezilla/scalezilla/cri"
)

// parseDeploymentFileSyntax only parse the provided file to check for syntax errors.
// It does not validate the structure of the content
func parseDeploymentFileSyntax(configFile string) error {
	parser := hclparse.NewParser()
	_, diagnostics := parser.ParseHCLFile(configFile)
	if diagnostics.HasErrors() {
		return diagnostics.Errs()[0]
	}

	return nil
}

// parseDeployment parses the deployment file spec to later create pods
func (c *Cluster) parseDeployment(data []byte) (cri.DeploymentSpec, error) {
	parser := hclparse.NewParser()
	file, diagnostics := parser.ParseHCL(data, "")
	if diagnostics.HasErrors() {
		return cri.DeploymentSpec{}, diagnostics.Errs()[0]
	}

	re := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]{5,62}$`)
	var spec cri.DeploymentSpec
	if diagnostics := gohcl.DecodeBody(file.Body, nil, &spec); diagnostics.HasErrors() {
		return cri.DeploymentSpec{}, diagnostics.Errs()[0]
	}

	if !re.MatchString(spec.Deployment.Name) {
		return cri.DeploymentSpec{}, ErrDeploymentNameInvalid
	}

	if !re.MatchString(spec.Deployment.Pod.Name) {
		return cri.DeploymentSpec{}, ErrDeploymentNameInvalid
	}

	if !re.MatchString(spec.Deployment.Pod.Container.Name) {
		return cri.DeploymentSpec{}, ErrDeploymentNameInvalid
	}

	if spec.Deployment.Namespace == "" {
		spec.Deployment.Namespace = "default"
	}

	if spec.Deployment.Kind == "" {
		spec.Deployment.Kind = "service"
	}

	if spec.Deployment.Pod.Container.Resources.CPU == 0 || spec.Deployment.Pod.Container.Resources.CPU <= 32 {
		spec.Deployment.Pod.Container.Resources.CPU = 128
	}

	if spec.Deployment.Pod.Container.Resources.Memory == 0 || spec.Deployment.Pod.Container.Resources.Memory <= 32 {
		spec.Deployment.Pod.Container.Resources.Memory = 128
	}

	return spec, nil
}
