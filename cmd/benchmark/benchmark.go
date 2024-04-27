package benchmark

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cecobask/timescale-coding-challenge/cmd"
	"github.com/cecobask/timescale-coding-challenge/internal/database"
	"github.com/cecobask/timescale-coding-challenge/pkg/log"
)

func NewCommand() *cobra.Command {
	var (
		configPath  string
		workerCount int
	)
	command := &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", cmd.CommandNameBenchmark),
		Short: "benchmark timescale database",
		PreRunE: func(c *cobra.Command, args []string) (err error) {
			configPath, err = c.Flags().GetString(cmd.FlagNameConfig)
			if err != nil {
				return err
			}
			workerCount, err = c.Flags().GetInt(cmd.FlagNameWorkers)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			db, err := database.New(c.Context())
			if err != nil {
				return err
			}
			defer db.Close()
			logger := log.FromContext(c.Context())
			logger.Info("benchmarking timescale database...", "config", configPath, "workers", workerCount)
			return nil
		},
	}
	command.Flags().StringP(cmd.FlagNameConfig, "c", "", "path to the config file")
	command.Flags().IntP(cmd.FlagNameWorkers, "w", 1, "number of workers to use")
	return command
}
