package member

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listTeamMemOptions struct {
	org   string
	team  string
	quiet bool
}

// NewListTeamMemCommand returns a new instance of the list team member command.
func NewListTeamMemCommand(c cli.Interface) *cobra.Command {
	opts := listTeamMemOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List members",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamMem(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display team member names")
	return cmd
}

func listTeamMem(c cli.Interface, cmd *cobra.Command, opts listTeamMemOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if opts.quiet {
		for _, member := range reply.Team.Members {
			c.Console().Println(member)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "MEMBER")
	for _, member := range reply.Team.Members {
		fmt.Fprintf(w, "%s\n", member)
	}
	w.Flush()
	return nil
}
