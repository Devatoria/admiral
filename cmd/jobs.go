package cmd

import (
	"fmt"
	"os"

	"github.com/Devatoria/admiral/jobs"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobRunCmd)
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Admiral jobs commander",
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available jobs",
	Run: func(cmd *cobra.Command, args []string) {
		for jobName := range jobs.Jobs {
			fmt.Println(jobName)
		}
	},
}

var jobRunCmd = &cobra.Command{
	Use:   "run [job]",
	Short: "Run provided job",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing job name, please use --help")
			return
		}

		jobName := args[0]
		if fn, ok := jobs.Jobs[jobName]; ok {
			err := fn(args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		} else {
			fmt.Println("Unknown job")
			return
		}
	},
}
