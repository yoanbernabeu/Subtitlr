/*
Copyright © 2023 Yoan Bernabeu <contact@yoandev.co>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a .env file with your OpenAI API key",
	Long:  `The 'configure' command creates a .env file with your OpenAI API key in your current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		/* Variables declaration */
		apiKey, _ := cmd.Flags().GetString("apiKey")

		//Check if the .env file already exists
		if _, err := os.Stat(".env"); os.IsNotExist(err) {
			//Create the .env file
			file, err := os.Create(".env")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			//Write the OpenAI API key in the .env file
			_, err = file.WriteString("OPENAI_API_KEY=" + apiKey)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("The .env file has been created in your current directory")

			//Display the success message
			fmt.Println("---------------------------------------")
			fmt.Println("You have entered the following values:")
			fmt.Println("apiKey:", apiKey)
			fmt.Println("---------------------------------------")

			return
		}

		//Display the error message
		fmt.Println("---------------------------------------")
		fmt.Println("The .env file already exists in your current directory")
		fmt.Println("---------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().String("apiKey", "", "OpenAI API key")
	configureCmd.MarkFlagRequired("apiKey")
}
