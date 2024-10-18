package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gabriel-panz/gomato/repository"
	"github.com/gabriel-panz/gomato/types"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var IsList bool = false

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a configuration for your timers",
	Long: `Create or list your pre-configured timers, you'll be guided through the creation by a sequence of prompts. For example:
create a new configuration:
	config
list all configurations:
	config -l
	config --list
`,
	Run: func(cmd *cobra.Command, args []string) {
		if IsList {
			listConfigs()
		} else {
			runConfigCommand(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&IsList, "list", "l", false, "lists all configurations")
}

func listConfigs() {
	r := repository.GetTimerRepo()
	res, err := r.GetAllConfigs()
	if err != nil {
		return
	}

	s := "\nYour configurations:\n"
	for _, c := range res {
		s += fmt.Sprintf("-- %s| work: %s rest: %s\n", c.Name, c.WorkTime, c.PauseTime)
	}

	fmt.Println(s)
}

func runConfigCommand(args []string) {
	conf := &types.TimerConfig{}
	var err error

	promptName := getPrompt(
		"Choose a name for this configuration (must be unique)",
		validateName)

	conf.Name, err = promptName.Run()
	if err != nil {
		log.Println(err)
		return
	}

	promptWorkTime := getPrompt(
		"How long would you like to work for? (use the format: 00h00m00s)",
		validateDurationInput)

	conf.WorkTime, err = runTimePrompt(promptWorkTime)
	if err != nil {
		log.Panicln(err)
		return
	}

	promptPauseTime := getPrompt(
		"How long would you like your break to last? (use the format: 00h00m00s)",
		validateDurationInput)

	conf.PauseTime, err = runTimePrompt(promptPauseTime)
	if err != nil {
		log.Panicln(err)
		return
	}

	sNotificationLevel := getSelect(
		"Choose a Notification Level",
		[]string{"None", "Audio Only", "Audio Visual"})

	i, _, err := sNotificationLevel.Run()
	if err != nil {
		log.Println(err)
		return
	}

	conf.NotificationLevel = types.NotificationLevel(i)
	r := repository.GetTimerRepo()
	err = r.InsertConfig(conf)
	if err != nil {
		log.Println(err)
		return
	}
}

func getPrompt(label string, validation func(s string) error) *promptui.Prompt {
	return &promptui.Prompt{
		Label:    label,
		Validate: validation,
	}
}

func getSelect(label string, items []string) *promptui.Select {
	return &promptui.Select{
		Label: label,
		Items: items,
	}
}

func runTimePrompt(p *promptui.Prompt) (time.Duration, error) {
	d, err := p.Run()
	if err != nil {
		return 0, err
	}

	r, err := time.ParseDuration(d)
	if err != nil {
		return 0, err
	}

	return r, nil
}

func validateDurationInput(input string) error {
	if len(input) <= 0 {
		return errors.New("please provide a non empty input")
	}

	_, err := time.ParseDuration(input)
	if err != nil {
		return err
	}

	return nil
}

func validateName(input string) error {
	if len(input) <= 0 {
		return errors.New("please provide a non empty input")
	}

	r := repository.GetTimerRepo()
	_, err := r.GetConfigByName(input)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			return nil
		}

		return nil
	} else {
		return errors.New("this name is already taken by another configuration")
	}
}
