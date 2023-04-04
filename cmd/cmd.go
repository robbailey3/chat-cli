package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/slog"
	"github.com/robbailey3/chat-cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "openai",
		Short: "A basic CLI for interacting with ChatGPT by openAI",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(
				ui.NewModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)

			if _, err := p.Run(); err != nil {
				slog.Error(err)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.openai-cli.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".chat-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		key := viper.Get("openAi.apiKey")
		err := os.Setenv("OPEN_AI_API_KEY", key.(string))
		if err != nil {
			slog.Error(err)
			return
		}
	} else {
		slog.Error(err)
	}
}
