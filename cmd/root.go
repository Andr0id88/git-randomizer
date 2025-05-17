package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"git-randomizer/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "gitr",
		Short: "Git randomizer: Rewrite git commit messages in outrageous personas",
		Long:  "Git randomizer: jazzes up your git life: commits, branches, and celebratory one-liners – all in character!",
	}
)

func Execute() { cobra.CheckErr(rootCmd.Execute()) }

func init() {
	home, _ := os.UserHomeDir()
	defaultCfg := filepath.Join(home, ".config", "git-randomizer", "git-randomizer.yaml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultCfg,
		"config file (default $HOME/.config/gitrandomizer/gitrandomizer.yaml)")
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(branchCmd)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	// sensible defaults (overridden by YAML)
	viper.SetDefault("default_character", "random")
	viper.SetDefault("default_group", "")
	viper.SetDefault("default_mood", "playful")
	viper.SetDefault("default_length", "medium")
	viper.SetDefault("confirm", true)

	viper.SetDefault("pass_secret", "gemini_api_key")

	viper.SetDefault("tagline_persona", "yoda")
	viper.SetDefault("tagline_enabled", true)

	viper.SetDefault("branch_persona", "random")
	viper.SetDefault("branch_persona_group", "")

	if err := viper.ReadInConfig(); err != nil {
		// first run – drop a commented sample file
		if err := config.CreateDefault(cfgFile); err != nil {
			fmt.Printf("config error: %v\n", err)
		} else {
			fmt.Printf("📝 Created default config at %s\n", cfgFile)
			_ = viper.ReadInConfig() // read the file we just wrote
		}
	}
}

