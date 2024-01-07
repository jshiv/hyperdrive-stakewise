/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/nodeset-org/hyperdrive-stakewise/hyperdrive"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Stops and restarts node containers.",
	Long: `Runs a "docker compose down" and "docker compose up -d".
	
The --clean option will run execute "docker compose down --remove-orphans" which will remove any containers not associated with your NodeSet-StakeWise configuration.
`,
	Run: func(cmd *cobra.Command, args []string) {
		clean, _ := cmd.Flags().GetBool("clean")
		color.HiWhite("Shutting down...")

		c, err := hyperdrive.LoadConfig()
		if err != nil {
			log.Fatal("Can't read config:", errors.Join(err, hyperdrive.ErrorCanNotFindConfigFile))
		}

		var removeOrphans string
		if clean {
			prompt := promptui.Select{
				Label: "This will remove any containers not associated with your NodeSet-StakeWise configurationAre, you sure you want to continue?",
				Items: []string{"n", "y"},
			}
			var err error
			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				log.Fatal(err)
			}
			if result == "y" {
				removeOrphans = "--remove-orphans"
			}
		}
		text := fmt.Sprintf("docker compose  -f compose.yaml -f compose.internal.yaml down %s", removeOrphans)
		log.Info(text)
		err = c.ExecCommand(text)
		if err != nil {
			log.Fatal(err)
		}
		color.HiWhite("Starting node...")

		var composeFile string
		if c.InternalClients {
			composeFile = "-f compose.yaml -f compose.internal.yaml"
		} else {
			composeFile = "-f compose.yaml"
		}
		text = fmt.Sprintf("docker compose %s pull", composeFile)
		log.Info(text)
		err = c.ExecCommand(text)
		if err != nil {
			log.Fatal(err)
		}

		text = fmt.Sprintf("docker compose %s up -d", composeFile)
		log.Info(text)
		err = c.ExecCommand(text)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().Bool("clean", false, "WARNING: Using the --clean option for stop will remove any containers not associated with your NodeSet-StakeWise configuration.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
