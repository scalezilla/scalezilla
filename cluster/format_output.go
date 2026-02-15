package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// printTableNodesList gives a table output
func printTableNodesList(body []byte) {
	var nodes []APINodesListResponse
	if err := json.Unmarshal(body, &nodes); err != nil {
		fmt.Println(err)
		return
	}

	// tabwriter aligns columns using tabs
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tADDRESS\tKIND\tLEADER\tNODEPOOL")

	for _, n := range nodes {
		leader := "false"
		if n.Leader {
			leader = "true"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			n.ID, n.Name, n.Address, strings.ToLower(n.Kind), leader, n.NodePool,
		)
	}
	w.Flush()
}
