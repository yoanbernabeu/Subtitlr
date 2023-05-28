/*
Copyright © 2023 Yoan Bernabeu <contact@yoandev.co>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Subtitlr",
	Short: "AI-assisted subtitle generation CLI for Youtube",
	Long: `This application, a subtitle generator for YouTube, utilizes OpenAI's Whisper API.

This tool leverages artificial intelligence to efficiently transcribe speech in YouTube videos into text, thereby generating accurate subtitles (in SRT format).

It's designed to improve the accessibility and convenience of video content, ensuring that no matter your language or hearing ability, you can fully engage with and comprehend the material.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Subtitlr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
