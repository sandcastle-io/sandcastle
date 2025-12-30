package flags

import "github.com/spf13/cobra"

// AddAllNamespacesFlagVar adds the --all-namespaces flag and binds it to target.
func AddAllNamespacesFlagVar(cmd *cobra.Command, target *bool) {
	cmd.Flags().BoolVarP(target, "all-namespaces", "A", false,
		"If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
}

// AddDryRunFlag adds dry-run flag to a command.
func AddDryRunFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"dry-run",
		"none",
		`Must be "none", "server", or "client". If client strategy, only print the object that would be sent, without sending it. If server strategy, submit server-side request without persisting the resource.`,
	)
}
