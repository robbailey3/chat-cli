package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/slog"
	"github.com/robbailey3/openai-cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "openai",
		Short: "A basic CLI for interacting with ChatGPT by openAI",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(ui.NewModel())

			if _, err := p.Run(); err != nil {
				slog.Error(err)
			}
			// client := openai.NewClient()
			//
			// completion, err := client.GetChatCompletion(context.Background(), openai.ChatCompletionRequest{
			//   Messages: []openai.ChatMessage{
			//     openai.ChatMessage{
			//       Role:    "user",
			//       Content: "Say how good I am at making ChatGPT work from my Golang CLI. Sing my praises!",
			//     },
			//   },
			//   Model:       "gpt-4",
			//   N:           1,
			//   Temperature: 1,
			//   TopP:        1,
			//   MaxTokens:   250,
			// })
			// if err != nil {
			//   slog.Error(err)
			//   return
			// }
			// slog.Info(completion)
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
		viper.SetConfigName(".openai-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file:", viper.ConfigFileUsed())
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
