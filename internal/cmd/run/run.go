package run

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/clientgetter"
	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd/flags"
	"github.com/mbobrovskyi/kube-sandcastle/pkg/api"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	runExample = templates.Examples(`
  		# Run local python file
		sandcastlectl run main.py

		# Run code directly from the console via stdin
		echo 'print("hello")' | sandcastlectl run -
	`)
)

type runOptions struct {
	code      string
	namespace string

	proxyClient rest.Interface

	streams genericiooptions.IOStreams
}

func newRunOptions(streams genericiooptions.IOStreams) *runOptions {
	return &runOptions{
		streams: streams,
	}
}

func NewRunCmd(clientGetter clientgetter.ClientGetter, streams genericiooptions.IOStreams) *cobra.Command {
	opts := newRunOptions(streams)

	cmd := &cobra.Command{
		Use:     "run [file]",
		Short:   "Run a code snippet in a secure sandboxed environment",
		Example: runExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if err := opts.complete(clientGetter, cmd, args); err != nil {
				return err
			}

			if err := opts.run(cmd.Context()); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AddDryRunFlag(cmd)

	return cmd
}

func (o *runOptions) complete(clientGetter clientgetter.ClientGetter, cmd *cobra.Command, args []string) error {
	var (
		code []byte
		err  error
	)

	if len(args) == 0 || args[0] == "-" {
		code, err = io.ReadAll(cmd.InOrStdin())
	} else {
		code, err = os.ReadFile(args[0])
	}
	if err != nil {
		return fmt.Errorf("failed to read code: %w", err)
	}
	o.code = string(code)

	o.namespace, _, err = clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("failed to get current namespace: %w", err)
	}

	o.proxyClient, err = clientGetter.ProxyClient()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %w", err)
	}

	return nil
}

func (o *runOptions) run(ctx context.Context) error {
	payload := api.ExecuteRequest{
		Code: o.code,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	res := o.proxyClient.Post().
		Namespace(o.namespace).
		Resource("services").
		Name("sandcastle-svc:http").
		SubResource("proxy").
		Suffix("execute").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Body(body).
		Do(ctx)
	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	raw, err := res.Raw()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	resp := &api.ExecuteResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	o.streams.Out.Write([]byte(resp.Stdout))
	o.streams.ErrOut.Write([]byte(resp.Stderr))

	return nil
}
