// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "codeprobot",
	Short: "CodePRobot is an automated PR assistant powered by ChatGPT and GitHub CLI",
	Long: `CodePRobot is a lightweight developer bot that watches file changes,
generates code using ChatGPT API, and automates GitHub PRs via gh CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 CodePRobot ready. Use --help for available commands.")
		if verbose {
			fmt.Println("🔍 Verbose mode enabled.")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
