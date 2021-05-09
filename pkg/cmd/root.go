package cmd

import (
	"github.com/spf13/cobra"
)

const (
	cliName = "ec2-connect"
)

func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   cliName,
		Short: "Connect, start, stop, list EC2 instances!",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
	}

	command.AddCommand(NewVersionCommand())
	command.AddCommand(NewListCommand())
	command.AddCommand(NewStartCommand())
	command.AddCommand(NewStopCommand())
	command.AddCommand(NewGroupCommand())
	command.AddCommand(NewConnectCommand())
	return command
}
