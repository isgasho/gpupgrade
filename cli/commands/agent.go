package commands

import (
	"github.com/greenplum-db/gp-common-go-libs/cluster"
	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/spf13/cobra"

	"github.com/greenplum-db/gpupgrade/agent"

	"github.com/greenplum-db/gpupgrade/utils"
	"github.com/greenplum-db/gpupgrade/utils/daemon"
	"github.com/greenplum-db/gpupgrade/utils/log"
)

func Agent() *cobra.Command {
	var logdir, statedir string
	var shouldDaemonize bool

	var cmd = &cobra.Command{
		Use:    "agent",
		Short:  "Start the Command Listener (blocks)",
		Long:   `Start the Command Listener (blocks)`,
		Hidden: true,
		Args:   cobra.MaximumNArgs(0), //no positional args allowed
		RunE: func(cmd *cobra.Command, args []string) error {
			gplog.InitializeLogging("gpupgrade agent", logdir)
			defer log.WritePanics()

			conf := agent.Config{
				Port:     6416,
				StateDir: statedir,
			}

			agentServer := agent.NewServer(&cluster.GPDBExecutor{}, conf)
			if shouldDaemonize {
				agentServer.MakeDaemon()
			}

			// blocking call
			agentServer.Start()

			return nil
		},
	}

	cmd.Flags().StringVar(&logdir, "log-directory", "", "command_listener log directory")
	cmd.Flags().StringVar(&statedir, "state-directory", utils.GetStateDir(), "Agent state directory")

	daemon.MakeDaemonizable(cmd, &shouldDaemonize)

	return cmd
}
