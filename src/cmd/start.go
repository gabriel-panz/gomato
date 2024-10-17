package cmd

import (
	"log"
	"time"

	"github.com/gabriel-panz/gomato/repository"
	"github.com/gabriel-panz/gomato/service"
	"github.com/gabriel-panz/gomato/types"
	"github.com/spf13/cobra"
)

var LocalWorkTime time.Duration = 0
var LocalPauseTime time.Duration = 0

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the pomodoro timer.",
	Long: `Start the pomodoro timer, can provide specific configuration. For example:
default configuration:
	start
using a custom configuration (see 'help config'):
	start <config name>
working for 20 minutes and pausing for 10 minutes:
	start -w 20m -p 10m
custom configuration, with modified pause for 10 hours and 32 minutes:
	start <config name> -p 10h32m
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runStartCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().DurationVarP(&LocalWorkTime, "work-time", "w", 0, "Amount of time of complete focus on work.\nFormat is 0h0m0s where h is hours, m is minutes, s is seconds")
	startCmd.Flags().DurationVarP(&LocalPauseTime, "pause-time", "p", 0, "Amount of time for relaxing.\nFormat is 0h0m0s where h is hours, m is minutes, s is seconds")
}

func runStartCommand(args []string) {
	t := service.CreateTimer()
	r := repository.GetTimerRepo()

	var (
		c   *types.TimerConfig
		err error
	)
	if len(args) > 0 {
		c, err = r.GetConfigByName(args[0])
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c, err = r.GetDefaultConfig()
		if err != nil {
			log.Println(err)
			return
		}
	}

	if LocalWorkTime > 0 {
		c.WorkTime = LocalWorkTime
	}

	if LocalPauseTime > 0 {
		c.PauseTime = LocalPauseTime
	}

	log.Printf("working for %v | pausing for %v\n", c.WorkTime, c.PauseTime)

	t.StartPomodoro(c)
}
