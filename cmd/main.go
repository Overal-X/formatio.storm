package main

import (
	"fmt"
	"os"

	storm "github.com/Overal-X/formatio.storm"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	version   = "dev" // default value
	commit    = "none"
	buildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "storm",
	Short: "Formatio Storm",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is the version of your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", version, commit, buildDate)
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Storm agent commands",
}

var agentRunWorkflowCmd = &cobra.Command{
	Use:  "run",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflowFile := args[0]
		inventoryFile, _ := cmd.Flags().GetString("inventory")
		format, _ := cmd.Flags().GetInt("format")

		agent := storm.NewAgent()
		err := agent.Run(
			agent.AgentWithFiles(workflowFile, inventoryFile),
			agent.AgentWithCallback(func(i interface{}) { fmt.Println(i) }, format),
		)
		if err != nil {
			os.Exit(1)
		}
	},
}

var agentInstallCmd = &cobra.Command{
	Use: "install",
	Run: func(cmd *cobra.Command, args []string) {
		inventoryFile, _ := cmd.Flags().GetString("inventory")
		installationMode, _ := cmd.Flags().GetString("mode")

		agent := storm.NewAgent()
		err := agent.Install(storm.InstallArgs{
			If:   inventoryFile,
			Mode: installationMode,
		})
		if err != nil {
			os.Exit(1)
		}
	},
}

var agentUninstallCmd = &cobra.Command{
	Use: "uninstall",
	Run: func(cmd *cobra.Command, args []string) {
		inventoryFile, _ := cmd.Flags().GetString("inventory")

		agent := storm.NewAgent()
		err := agent.Uninstall(storm.UninstallArgs{If: inventoryFile})
		if err != nil {
			os.Exit(1)
		}
	},
}

var runWorkflowCmd = &cobra.Command{
	Use:  "run",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflowFile := args[0]
		trashWorkflow, _ := cmd.Flags().GetBool("trash-workflow")
		directory, _ := cmd.Flags().GetString("directory")
		format, _ := cmd.Flags().GetInt("format")

		if trashWorkflow {
			defer os.Remove(workflowFile)
		}

		workflow := storm.NewWorkflow()

		wc, err := workflow.Load(workflowFile)
		if err != nil {
			os.Exit(1)
		}

		wc.Directory = lo.Ternary(wc.Directory == "" && directory != "", directory, wc.Directory)

		err = workflow.Run(
			workflow.WorkflowWithConfig(*wc),
			workflow.WorkflowWithCallback(func(i interface{}) { fmt.Println(i) }, format))
		if err != nil {
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.AddCommand(versionCmd)

	agentInstallCmd.Flags().StringP("inventory", "i", "./inventory.yaml", "formatio storm inventory")
	agentInstallCmd.Flags().StringP("mode", "m", "prod", "formatio storm installation type (prod or dev)")
	agentCmd.AddCommand(agentInstallCmd)

	agentUninstallCmd.Flags().StringP("inventory", "i", "./inventory.yaml", "formatio storm inventory")
	agentCmd.AddCommand(agentUninstallCmd)

	agentRunWorkflowCmd.Flags().StringP("inventory", "i", "./inventory.yaml", "formatio storm inventory")
	agentRunWorkflowCmd.Flags().IntP("format", "f", 1, "available options are; 1 => plain, 2 => struct, 3 => json")
	agentCmd.AddCommand(agentRunWorkflowCmd)

	runWorkflowCmd.Flags().BoolP("trash-workflow", "t", true, "remove workflow file if the workflow is complete")
	runWorkflowCmd.Flags().StringP("directory", "d", ".", "directory to run the workflow from")
	runWorkflowCmd.Flags().IntP("format", "f", 1, "available options are; 1 => plain, 2 => struct, 3 => json")
	rootCmd.AddCommand(runWorkflowCmd)

	rootCmd.AddCommand(agentCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // TODO: use logger
		os.Exit(1)
	}
}
