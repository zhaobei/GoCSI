package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

var nodeGetInfoCmd = &cobra.Command{
	Use:     "get-info",
	Aliases: []string{"info"},
	Short:   `invokes the rpc "NodeGetInfo"`,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithTimeout(root.ctx, root.timeout)
		defer cancel()

		rep, err := node.client.NodeGetInfo(ctx, &csi.NodeGetInfoRequest{})
		if err != nil {
			return err
		}

		return root.tpl.Execute(os.Stdout, rep)
	},
}

func init() {
	nodeCmd.AddCommand(nodeGetInfoCmd)
}
