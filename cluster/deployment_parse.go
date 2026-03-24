package cluster

import "github.com/hashicorp/hcl/v2/hclparse"

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
