package cmd

import (
	"fmt"
	"os"

	"github.com/chelnak/godf/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "godf",
	Short: "I am a TUI app that watches Azure Data Factory pipeline runs.",
	Long:  "I am a TUI app that watches Azure Data Factory pipeline runs.",

	Run: func(cmd *cobra.Command, args []string) {
		ui.Draw()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.catz.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".godf")
	}

	viper.AutomaticEnv()

	viper.SetDefault("RefreshIntervalSeconds", "30")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Could not find config file .godf in your home directory!")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
