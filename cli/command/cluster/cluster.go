package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
)

type clusterOpts struct {
	managers      int
	workers       int
	name          string
	tag           string
	provider      string
	region				string
	domain				string
	registration  string
	notifications bool
	log						int
}

const (
	DefaultLocalClusterID = "f573e897-7aa0-4516-a195-42ee91039e97"
)

var (
	opts = &clusterOpts{
		managers:      3,
		workers:       2,
		name:          "",
		tag:           "latest",
		provider:      "local",
		region:				 "",
		domain:				 "",
		registration:  configuration.RegistrationDefault,
		notifications: true,
		log:			 		 4,
	}
)

// NewClusterCommand returns a new instance of the cluster command.
func NewClusterCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Cluster management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
	cmd.AddCommand(NewNodeCommand(c))
	return cmd
}

func queryCluster(c cli.Interface, args []string, env map[string]string) error {
	if err := check(opts.provider); err != nil {
		return err
	}
	err := Run(c, args, env)
	if err != nil {
		// TODO: the local cluster is the only one that can be managed this release
		c.Console().Println(DefaultLocalClusterID)
	}
	return err
}

// Map cli cluster flags to target bootstrap cluster command flags,
// append to and return args array
func reflag(cmd *cobra.Command, flags map[string]string, args []string) []string {
	// transform src flags to target flags and add flag and value to cargs
	for s, t := range flags {
		if cmd.Flag(s).Changed {
			args = append(args, t, cmd.Flag(s).Value.String())
		}
	}
	return args
}
