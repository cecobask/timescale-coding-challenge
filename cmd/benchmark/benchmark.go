package benchmark

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cecobask/timescale-coding-challenge/cmd"
	"github.com/cecobask/timescale-coding-challenge/internal/database"
	"github.com/cecobask/timescale-coding-challenge/internal/orchestrator"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", cmd.CommandNameBenchmark),
		Short: "benchmark timescale database",
		PreRunE: func(c *cobra.Command, args []string) (err error) {
			fileInfo, err := os.Stdin.Stat()
			if err != nil {
				return err
			}
			configProvided := c.Flags().Lookup(cmd.FlagNameConfig).Changed || (fileInfo.Mode()&os.ModeNamedPipe != 0)
			if !configProvided {
				return fmt.Errorf("config file or standard input must be provided")
			}
			return nil
		},
		RunE: run,
	}
	command.Flags().StringP(cmd.FlagNameConfig, "c", "", "path to the config file")
	command.Flags().IntP(cmd.FlagNameWorkers, "w", 1, "number of workers to use")
	return command
}

func run(c *cobra.Command, _ []string) error {
	configPath, err := c.Flags().GetString(cmd.FlagNameConfig)
	if err != nil {
		return err
	}
	workerCount, err := c.Flags().GetInt(cmd.FlagNameWorkers)
	if err != nil {
		return err
	}
	db, err := database.New(c.Context())
	if err != nil {
		return err
	}
	defer db.Close()
	o := orchestrator.New(workerCount, db)
	if err != nil {
		return err
	}
	if err = o.LoadConfig(configPath); err != nil {
		return err
	}
	benchmarks, err := o.Orchestrate(c.Context())
	if err != nil {
		return err
	}
	fmt.Println("Benchmarks count:", len(benchmarks))
	return nil
}
