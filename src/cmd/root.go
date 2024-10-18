package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gomato",
	Short: "An easy to use Pomodoro for the CLI.",
	Long: `Gomato is an easy to use Pomodoro/Flowmodoro for the CLI.
You can use the default configurations, change the timers for every work session
or create a configuration to save your work/break timer preferences. Get started with:

  gomato start
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
