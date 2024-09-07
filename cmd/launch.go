/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package cmd

import (
	// "fmt"

	"fmt"

	"github.com/konstructio/kubefirst/internal/launch"
	// "github.com/konstructio/kubefirst/internal/progress"
	"github.com/spf13/cobra"
)

// additionalHelmFlags can optionally pass user-supplied flags to helm
var additionalHelmFlags []string

func LaunchCommand() *cobra.Command {
	launchCommand := &cobra.Command{
		Use:   "launch",
		Short: "create a local k3d cluster and launch the kubefirst console and api in it",
		Long:  "create a local k3d cluster and launch the kubefirst console and api in it",
	}

	// wire up new commands
	launchCommand.AddCommand(launchUp(), launchDown(), launchCluster())

	return launchCommand
}

// launchUp creates a new k3d cluster with Kubefirst console and API
func launchUp() *cobra.Command {
	launchUpCmd := &cobra.Command{
		Use:              "up",
		Short:            "launch new console and api instance",
		TraverseChildren: true,
		// PreRun:           common.CheckDocker, // TODO: check runtimes when we can support more runtimes
		Run: func(cmd *cobra.Command, args []string) {
			launch.Up(additionalHelmFlags, false, true)
		},
	}

	launchUpCmd.Flags().StringSliceVar(&additionalHelmFlags, "helm-flag", []string{}, "additional helm flag to pass to the launch up command - can be used any number of times")

	return launchUpCmd
}

// launchDown destroys a k3d cluster for Kubefirst console and API
func launchDown() *cobra.Command {
	launchDownCmd := &cobra.Command{
		Use:              "down",
		Short:            "remove console and api instance",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			launch.Down(false)
		},
	}

	return launchDownCmd
}

// launchCluster
func launchCluster() *cobra.Command {
	launchClusterCmd := &cobra.Command{
		Use:              "cluster",
		Short:            "interact with clusters created by the kubefirst console",
		TraverseChildren: true,
	}

	launchClusterCmd.AddCommand(launchListClusters(), launchDeleteCluster())

	return launchClusterCmd
}

// launchListClusters makes a request to the console API to list created clusters
func launchListClusters() *cobra.Command {
	launchListClustersCmd := &cobra.Command{
		Use:              "list",
		Short:            "list clusters created by the kubefirst console",
		TraverseChildren: true,
		// PreRun:           common.CheckDocker,
		Run: func(cmd *cobra.Command, args []string) {
			launch.ListClusters()
		},
	}

	return launchListClustersCmd
}

// launchDeleteCluster makes a request to the console API to delete a single cluster
func launchDeleteCluster() *cobra.Command {
	launchDeleteClusterCmd := &cobra.Command{
		Use:              "delete",
		Short:            "delete a cluster created by the kubefirst console",
		TraverseChildren: true,
		// PreRun:           common.CheckDocker,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("you must provide a cluster name as the first argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] is the cluster name
			launch.DeleteCluster(args[0])
		},
	}
	launchDeleteClusterCmd.Flags().BoolVar(&ciFlag, "ci", false, "if running kubefirst in ci, set this flag to disable interactive features")
	return launchDeleteClusterCmd
}
