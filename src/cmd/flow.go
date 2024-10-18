package cmd

import (
	"github.com/gabriel-panz/gomato/service"
	"github.com/spf13/cobra"
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Starts a Pomodoro-esque timer that focuses on preserving the flow state.",
	Long:  `Starts a timer that keeps counting up, then once you begin a break it divides the spent time by 5 and starts another timer.`,
	Run: func(cmd *cobra.Command, args []string) {
		runFlowTimer()
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
}

func runFlowTimer() {
	t := service.CreateTimer()
	t.StartFlow()
}
