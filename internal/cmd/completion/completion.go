package completion

import (
	"strings"

	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/clientgetter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const completionLimit = 100

func NamespaceNameFunc(clientGetter clientgetter.ClientGetter) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clientSet, err := clientGetter.K8sClientSet()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}

		list, err := clientSet.CoreV1().Namespaces().List(cmd.Context(), metav1.ListOptions{Limit: completionLimit})
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}

		validArgs := make([]string, len(list.Items))
		for i, wl := range list.Items {
			validArgs[i] = wl.Name
		}

		return validArgs, cobra.ShellCompDirectiveNoFileComp
	}
}

func ContextsFunc(clientGetter clientgetter.ClientGetter) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := clientGetter.ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var validArgs []string
		for name := range config.Contexts {
			if strings.HasPrefix(name, toComplete) {
				validArgs = append(validArgs, name)
			}
		}
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	}
}

func ClustersFunc(clientGetter clientgetter.ClientGetter) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := clientGetter.ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var validArgs []string
		for name := range config.Clusters {
			if strings.HasPrefix(name, toComplete) {
				validArgs = append(validArgs, name)
			}
		}
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	}
}

func UsersFunc(clientGetter clientgetter.ClientGetter) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := clientGetter.ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var validArgs []string
		for name := range config.AuthInfos {
			if strings.HasPrefix(name, toComplete) {
				validArgs = append(validArgs, name)
			}
		}
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	}
}
