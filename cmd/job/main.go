package main

import (
	"encoding/json"
	"fmt"
	"os"

	"charlotte/internal/job"
	jobrun "charlotte/internal/jobrun"
	localruntime "charlotte/internal/runtime/local"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "job",
		Short: "Run job locally",
	}

	runLocalCmd := &cobra.Command{
		Use:   "run-local",
		Short: "Runs YAML job file",
		RunE:  runJobHandler,
	}

	runLocalCmd.Flags().StringP("job", "j", "", "Path to filename with a job")
	_ = runLocalCmd.MarkFlagRequired("job")
	runLocalCmd.Flags().StringP("inputs", "i", "", "Path to file containing input values")
	_ = runLocalCmd.MarkFlagRequired("inputs")
	runLocalCmd.Flags().BoolP("quiet", "q", false, "Do not print step stdout and stderr")
	runLocalCmd.Flags().StringP("result", "r", "", "Path to write JSON result")

	rootCmd.AddCommand(runLocalCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version",
		Run:   versionHandler,
	}
	rootCmd.AddCommand(versionCmd)

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionHandler(cmd *cobra.Command, args []string) {
	fmt.Fprintln(os.Stdout, VERSION)
}

func runJobHandler(cmd *cobra.Command, args []string) error {
	jobFile, _ := cmd.Flags().GetString("job")
	j, err := job.NewFromFile(jobFile)
	if err != nil {
		return fmt.Errorf("error parsing job from file %s: %w", jobFile, err)
	}

	inputsFile, _ := cmd.Flags().GetString("inputs")
	b, err := os.ReadFile(inputsFile)
	if err != nil {
		return fmt.Errorf("error reading inputs file %s: %w", inputsFile, err)
	}

	var jobRunInputs jobrun.JobRunInputs
	err = json.Unmarshal(b, &jobRunInputs)
	if err != nil {
		return fmt.Errorf("error unmarshalling inputs file %s: %w", inputsFile, err)
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	runenv := localruntime.NewLocalRuntime(quiet)
	jobRunResult := j.Run(runenv, &jobRunInputs)
	if !quiet && !jobRunResult.Success {
		fmt.Fprintf(os.Stderr, "Job %s failed locally with: %s\n", jobFile, jobRunResult.Error.Error())
	}

	resultJson, err := json.Marshal(jobRunResult)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stdout, "Marshalling run result to JSON failed\n")
		}
		return err
	}

	outputFile, _ := cmd.Flags().GetString("result")
	if outputFile != "" {
		err = os.WriteFile(outputFile, resultJson, 0600)
		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stdout, "Writing JSON run result to file failed\n")
			}
			return err
		}
	} else {
		fmt.Fprintln(os.Stdout, string(resultJson))
	}

	if !jobRunResult.Success {
		return fmt.Errorf("job execution failed")
	}

	return nil
}
