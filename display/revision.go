package display

import (
	"fmt"

	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/revision"
)

// Revision displays a revision
func RevisionDiff(diff *revision.Diff, cloudService string, root *node.Node, verbose bool, showProperties bool) {
	fromRevision := "repository creation"
	if diff.From.ID != "" {
		fromRevision = diff.From.ID[:7] + " on " + diff.From.Time.Format("Monday January 2, 15:04")
	}

	if showProperties {
		if diff.GraphDiff.HasDiff() {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Time.Format("Monday January 2, 15:04"))
			FullDiff(diff.GraphDiff, root, cloudService)
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Time.Format("Monday January 2, 15:04"))
			fmt.Println("No changes.")
		}
	} else {
		if diff.GraphDiff.HasResourceDiff() {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Time.Format("Monday January 2, 15:04"))
			ResourceDiff(diff.GraphDiff, root)
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Time.Format("Monday January 2, 15:04"))
			fmt.Println("No resource changes.")
		}
	}
}
