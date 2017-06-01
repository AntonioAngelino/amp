package cluster

import (
	"strconv"

	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
)

// NewCreateCommand returns a new instance of the create command for bootstrapping an cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Intial number of manager nodes")
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	flags.StringVarP(&opts.tag, "tag", "t", "0.11.0", "Specify tag for cluster images, use 'local' for development")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider, options are local, docker and aws")
	flags.StringVar(&opts.region, "region", "", "Specify region for deployment on selected cloud provider")
	flags.StringVar(&opts.domain, "domain", "", "Specify the ssh key for deployment on selected cloud provider")
	flags.StringVarP(&opts.registration, "registration", "r", configuration.RegistrationNone, "Specify the registration policy, options are 'none' or 'email'")
	flags.BoolVarP(&opts.notifications, "notifications", "n", false, "Enable/disable server notifications")
	flags.IntVarP(&opts.log, "log", "l", 4, "Logging level. 0 is least verbose, Max is 5")
	return cmd
}

// Map cli cluster flags to target bootstrap cluster command flags and update the cluster
func create(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		"workers":       "-w",
		"managers":      "-m",
		"name":          "-i",
		"tag":           "-T",
		"provider":      "-p",
		"region":				 "-g",
		"domain":				 "-D",
		"registration":	 "-r",
		"notifications": "-n",
		"log":					 "-l",
	}

	// the following ensures that flags are added before the final command arg
	// TODO: refactor reflag to handle this
	args := []string{"bin/deploy"}
	args = reflag(cmd, m, args)
	env := map[string]string{"TAG": opts.tag, "REGION": opts.region, "DOMAIN": opts.domain, "REGISTRATION": opts.registration, "NOTIFICATIONS": strconv.FormatBool(opts.notifications)}
	return queryCluster(c, args, env)
}
