package cluster

import (
	"encoding/json"
	"errors"
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
	_, _ = fmt.Fprintln(w, "ID\tNAME\tADDRESS\tKIND\tLEADER\tNODEPOOL")

	for _, n := range nodes {
		leader := "false"
		if n.Leader {
			leader = "true"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			n.ID, n.Name, n.Address, strings.ToLower(n.Kind), leader, n.NodePool,
		)
	}
	_ = w.Flush()
}

// decodeError decode the current error to return it to the main func
func decodeError(body []byte) error {
	var resp respError
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	return errors.New(resp.Error)
}

// printTableNodesList gives a table output
func printTablePodsList(body []byte) {
	var z []APIPodsListResponse
	if err := json.Unmarshal(body, &z); err != nil {
		fmt.Println(err)
		return
	}

	// tabwriter aligns columns using tabs
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tSTATUS\tAGE")

	for _, n := range z {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n",
			n.ID, n.Status, n.CreatedAt,
		)
	}
	_ = w.Flush()
}

// printTablePodsDelete gives a table output
func printTablePodsDelete(body []byte) {
	var z APIPodsDeleteResponse
	if err := json.Unmarshal(body, &z); err != nil {
		fmt.Println(err)
		return
	}
	for _, x := range z.Pods {
		fmt.Println(x)
	}
}
