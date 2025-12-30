package dryrun

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Strategy int

const (
	// None indicates the client will make all mutating calls
	None Strategy = iota

	// Client or client-side dry-run, indicates the client will prevent
	// making mutating calls such as CREATE, PATCH, and DELETE
	Client

	// Server or server-side dry-run, indicates the client will send
	// mutating calls to the APIServer with the dry-run parameter to prevent
	// persisting changes.
	//
	// Note that clients sending server-side dry-run calls should verify that
	// the APIServer and the resource supports server-side dry-run, and otherwise
	// clients should fail early.
	//
	// If a client sends a server-side dry-run call to an APIServer that doesn't
	// support server-side dry-run, then the APIServer will persist changes inadvertently.
	Server
)

func GetStrategy(cmd *cobra.Command) (Strategy, error) {
	dryRunFlag, err := cmd.Flags().GetString("dry-run")
	if err != nil {
		return None, err
	}
	switch dryRunFlag {
	case "client":
		return Client, nil
	case "server":
		return Server, nil
	case "none":
		return None, nil
	default:
		return None, fmt.Errorf(`invalid dry-run value (%v). Must be "none", "server", or "client"`, dryRunFlag)
	}
}

// PrintFlagsWithStrategy sets a success message at print time for the dry run strategy
func PrintFlagsWithStrategy(printFlags *genericclioptions.PrintFlags, strategy Strategy) error {
	switch strategy {
	case Client:
		if err := printFlags.Complete("%s (client dry run)"); err != nil {
			return err
		}
	case Server:
		if err := printFlags.Complete("%s (server dry run)"); err != nil {
			return err
		}
	}
	return nil
}
