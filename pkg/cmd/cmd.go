package cmd

import (
	"errors"
	"fmt"
	"github.com/alicek106/go-ec2-ssh-autoconnect/pkg/aws"
	"github.com/alicek106/go-ec2-ssh-autoconnect/pkg/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

var (
	key string
	version      = "v0.7" // TODO : version should be injected in build time
	buildDate    = "1970-01-01T00:00:00Z"
	gitCommit    = ""
)

func NewListCommand() *cobra.Command {
	versionCmd := cobra.Command{
		Use:   "list",
		Short: "List EC2 instances",
		Example: `  # List all EC2 instances.
  ec2-connect list`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()
			svc.ListInstances()
		},
	}
	return &versionCmd
}

func NewStartCommand() *cobra.Command {
	startCmd := cobra.Command{
		Use:   "start",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Require a instance name.")
			}
			return nil
		},
		Short: "Start EC2 instances",
		Example: `  # Start a EC2 instances.
  ec2-connect start myserver`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()
			instanceName := args[0]
			instanceIDs := svc.GetInstanceIDs([]string{instanceName})
			for index, instanceName := range []string{instanceName} {
				log.Printf("Starting EC2 instance : %s (instance ID: %s)", instanceName, *instanceIDs[index])
			}
			svc.StartInstances(instanceIDs)
		},
	}
	return &startCmd
}

func NewStopCommand() *cobra.Command {
	stopCmd := cobra.Command{
		Use:   "stop",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Require a instance name.")
			}
			return nil
		},
		Short: "Stop EC2 instances",
		Example: `  # Stop a EC2 instances.
  ec2-connect stop myserver`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()
			instanceIDs := svc.GetInstanceIDs([]string{args[0]})
			for index, instanceName := range []string{args[0]} {
				log.Printf("Stoping EC2 instance : %s (instance ID: %s)", instanceName, *instanceIDs[index])
			}
			svc.StopInstances(instanceIDs)
		},
	}
	return &stopCmd
}

func NewGroupCommand() *cobra.Command {
	groupCmd := cobra.Command{
		Use:   "group",
		Short: "Manage EC2 instance as a group",
		Example: `  # Start a group of instances
  ec2-connect group start mygroup

  # Stop EC2 instance as a group
  ec2-connect group stop mygroup`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
			os.Exit(1)
		},
	}

	groupCmd.AddCommand(newGroupStartCommand())
	groupCmd.AddCommand(newGroupStopCommand())
	return &groupCmd
}

func newGroupStartCommand() *cobra.Command {
	groupStartCmd := cobra.Command{
		Use:   "start",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Require a group name.")
			}
			return nil
		},
		Short: "Start EC2 instance as a group",
		Example: `  # Start a group of instances
  ec2-connect group start mygroup`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()
			instanceNames := config.GetEnvparser().GetGroupInstanceNames(args[0])
			instanceIDs := svc.GetInstanceIDs(instanceNames)
			for index, instanceName := range instanceNames {
				log.Printf("Starting EC2 instance : %s (instance ID: %s)", instanceName, *instanceIDs[index])
			}
			svc.StartInstances(instanceIDs)
		},
	}
	return &groupStartCmd
}

func newGroupStopCommand() *cobra.Command {
	groupStopCmd := cobra.Command{
		Use:   "stop",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Require a group name.")
			}
			return nil
		},
		Short: "Stop EC2 instance as a group",
		Example: `  # Stop a group of instances
  ec2-connect group stop mygroup`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()
			instanceNames := config.GetEnvparser().GetGroupInstanceNames(args[0])
			instanceIDs := svc.GetInstanceIDs(instanceNames)
			for index, instanceName := range instanceNames {
				log.Printf("Stoping EC2 instance : %s (instance ID: %s)", instanceName, *instanceIDs[index])
			}
			svc.StopInstances(instanceIDs)
		},
	}
	return &groupStopCmd
}

func NewConnectCommand() *cobra.Command {
	connectCmd := cobra.Command{
		Use:   "connect",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Require a instance name.")
			}
			return nil
		},
		Short: "Connect to a EC2 instance using SSH",
		Example: `  # Connect to a EC2 instance
  ec2-connect connect myserver

  # Connect to a EC2 instance using user-defined key in configuration file.
  ec2-connect connect myserver --key=mykey`,
		Run: func(cmd *cobra.Command, args []string) {
			svc := aws.GetEC2Service()

			if cmd.Flags().Changed("key") {
				keyName, err := cmd.Flags().GetString("key")
				key = config.GetEnvparser().GetCustomKey(keyName)
				if err != nil {
					log.Fatalf("Something wrong to get key flag: %s", err)
				}
			}
			instanceName := args[0]
			instanceID := svc.GetInstanceIDs([]string{instanceName})
			if svc.GetInstanceStatus(instanceID) == "running" {
				log.Println("Instance in active.")
			} else {
				svc.StartInstances(instanceID)
				svc.WaitUntilActive(instanceID, []string{instanceName})
			}

			instanceIP := svc.GetInstancePublicIP(instanceName)
			sshUserName := svc.GetUsernamePerOS(instanceName)
			execCmd := exec.Command("ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("%s@%s", sshUserName, instanceIP), fmt.Sprintf("-i%s", key))
			execCmd.Stdin = os.Stdin
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			err := execCmd.Run()
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
		},
	}

	defaultKey := config.GetEnvparser().GetDefaultKey()
	connectCmd.Flags().StringVar(&key, "key", defaultKey, "SSH key to connect EC2 instance")
	return &connectCmd
}

func NewVersionCommand() *cobra.Command {
	versionCmd := cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Example: `  # Print version of binary.
  ec2-connect version
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("  Version: %s\n", version)
			fmt.Printf("  Build Date: %s\n", buildDate)
			fmt.Printf("  Git Commit: %s\n", gitCommit)
		},
	}
	return &versionCmd
}
