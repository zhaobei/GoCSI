package cmd

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

var nodePublishVolume struct {
	targetPath        string
	stagingTargetPath string
	pubCtx            mapOfStringArg
	volCtx            mapOfStringArg
	attribs           mapOfStringArg
	readOnly          bool
	caps              volumeCapabilitySliceArg
}

var nodePublishVolumeCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"mnt", "mount"},
	Short:   `invokes the rpc "NodePublishVolume"`,
	Example: `
USAGE

    csc node publish [flags] VOLUME_ID [VOLUME_ID...]
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		req := csi.NodePublishVolumeRequest{
			StagingTargetPath: nodePublishVolume.stagingTargetPath,
			TargetPath:        nodePublishVolume.targetPath,
			PublishContext:    nodePublishVolume.pubCtx.data,
			Readonly:          nodePublishVolume.readOnly,
			Secrets:           root.secrets,
			VolumeContext:     nodePublishVolume.volCtx.data,
		}

		if len(nodePublishVolume.caps.data) > 0 {
			req.VolumeCapability = nodePublishVolume.caps.data[0]
		}

		for i := range args {
			ctx, cancel := context.WithTimeout(root.ctx, root.timeout)
			defer cancel()

			// Set the volume ID for the current request.
			req.VolumeId = args[i]

			log.WithField("request", req).Debug("mounting volume")
			_, err := node.client.NodePublishVolume(ctx, &req)
			if err != nil {
				return err
			}

			fmt.Println(args[i])
		}

		return nil
	},
}

func init() {
	nodeCmd.AddCommand(nodePublishVolumeCmd)

	flagStagingTargetPath(
		nodePublishVolumeCmd.Flags(), &nodePublishVolume.stagingTargetPath)

	flagTargetPath(
		nodePublishVolumeCmd.Flags(), &nodePublishVolume.targetPath)

	flagReadOnly(
		nodePublishVolumeCmd.Flags(), &nodePublishVolume.readOnly)

	flagVolumeContext(nodePublishVolumeCmd.Flags(), &nodePublishVolume.volCtx)

	flagPublishContext(nodePublishVolumeCmd.Flags(), &nodePublishVolume.pubCtx)

	flagVolumeCapability(
		nodePublishVolumeCmd.Flags(), &nodePublishVolume.caps)

	flagWithRequiresVolContext(
		nodePublishVolumeCmd.Flags(), &root.withRequiresVolContext, false)

	flagWithRequiresPubContext(
		nodePublishVolumeCmd.Flags(), &root.withRequiresPubContext, false)

	flagWithRequiresCreds(
		nodePublishVolumeCmd.Flags(), &root.withRequiresCreds, "")
}
