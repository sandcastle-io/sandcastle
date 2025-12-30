package cmd

import (
	"os"

	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/clientgetter"
	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/completion"
	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/run"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/utils/clock"
)

type SandcastleOptions struct {
	Clock       clock.Clock
	ConfigFlags *genericclioptions.ConfigFlags

	genericiooptions.IOStreams
}

func defaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDiscoveryQPS(50.0)
}

func NewDefaultSandcastlectlCmd() *cobra.Command {
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	return NewSandcastlectlCmd(SandcastleOptions{
		ConfigFlags: defaultConfigFlags().WithWarningPrinter(ioStreams),
		IOStreams:   ioStreams,
		Clock:       clock.RealClock{},
	})
}

func NewSandcastlectlCmd(o SandcastleOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sandcastlectl",
		Short: "High-Performance, Secure, and Kueue-Native Runtime for AI Agents",
	}

	flags := cmd.PersistentFlags()

	configFlags := o.ConfigFlags
	if configFlags == nil {
		configFlags = defaultConfigFlags().WithWarningPrinter(o.IOStreams)
	}
	configFlags.AddFlags(flags)

	clientGetter := clientgetter.New(configFlags)

	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("namespace", completion.NamespaceNameFunc(clientGetter)))
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("context", completion.ContextsFunc(clientGetter)))
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("cluster", completion.ClustersFunc(clientGetter)))
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("user", completion.UsersFunc(clientGetter)))

	cmd.AddCommand(run.NewRunCmd(clientGetter, o.IOStreams))

	return cmd
}
