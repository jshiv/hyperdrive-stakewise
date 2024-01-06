/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
	"github.com/nodeset-org/hyperdrive-stakewise/hyperdrive"
	"github.com/nodeset-org/hyperdrive-stakewise/local"
)

// var Local embed.FS

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalizes the ~/.node-data/ directory with nodeset.env, compose.yaml and the ec and cc docker files.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("{::} Welcome to the NodeSet config script for StakeWise {::}")

		remove, _ := cmd.Flags().GetBool("remove")
		if remove {
			err := os.RemoveAll(dataDir)
			if err != nil {
				log.Fatal(err)
			}
		}
		network, _ := cmd.Flags().GetString("network")
		ecName, _ := cmd.Flags().GetString("ecname")
		checkpoint, _ := cmd.Flags().GetBool("checkpoint")

		if ecName == "" {
			prompt := promptui.Select{
				Label: "Select Execution Client",
				Items: []string{"geth", "nethermind"},
			}
			var err error
			_, ecName, err = prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				log.Fatal(err)
			}
		}

		ccName, _ := cmd.Flags().GetString("ccname")
		if ccName == "" {
			prompt := promptui.Select{
				Label: "Select Concensus Client",
				Items: []string{"nimbus", "teku"},
			}
			var err error
			_, ccName, err = prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				log.Fatal(err)
			}
		}

		var err error
		var c hyperdrive.Config
		if network == "" {
			prompt := promptui.Select{
				Label: "Select Network",
				Items: []string{"NodeSet Test Vault (holesky)", "Gravita (mainnet)", "NodeSet Dev Vault (holskey-dev)"},
			}
			var err error
			_, network, err = prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				log.Fatal(err)
			}
			switch network {
			case "NodeSet Test Vault (holesky)":
				c = hyperdrive.Holskey
			case "NodeSet Dev Vault (holskey-dev)":
				c = hyperdrive.HoleskyDev
			case "Gravita (mainnet)":
				c = hyperdrive.Gravita

			default:
				log.Fatalf("network %s is not avaliable, please choose holskey, holskey-dev or Gravita", network)
			}
		}

		if checkpoint {
			prompt := promptui.Prompt{
				Label:   "Provide Checkpoint Sync URL",
				Default: "https://checkpoint-sync.holesky.ethpandaops.io",
			}
			var err error
			checkpointURL, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				log.Fatal(err)
			}
			c.Checkpoint = checkpoint
			c.CheckpointURL = checkpointURL
		}

		dataDir := viper.GetString("DATA_DIR")
		c.DataDir = dataDir
		log.Infof("Writing config to data path: %s", dataDir)
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			log.Error(err)
		}

		c.SetViper()
		//Ensure that nodeset.env contains the correct ECNAME and CCNAME
		viper.Set("ECNAME", ecName)
		viper.Set("CCNAME", ccName)

		err = c.WriteConfig()
		if err != nil {
			log.Fatal(err)
		}

		//Write the compose file
		err = os.WriteFile(filepath.Join(dataDir, "compose.yaml"), local.Compose, 0766)
		if err != nil {
			log.Fatal(err)
		}

		//Select EL client
		ecCompose, err := local.Clients.ReadFile(fmt.Sprintf("clients/%s.yaml", ecName))
		if err != nil {
			log.Error(err)
		}
		err = os.WriteFile(filepath.Join(dataDir, fmt.Sprintf("%s.yaml", ecName)), ecCompose, 0766)
		if err != nil {
			log.Fatal(err)
		}

		//Select CC client
		ccCompose, err := local.Clients.ReadFile(fmt.Sprintf("clients/%s.yaml", ccName))
		if err != nil {
			log.Error(err)
		}
		err = os.WriteFile(filepath.Join(dataDir, fmt.Sprintf("%s.yaml", ccName)), ccCompose, 0766)
		if err != nil {
			log.Fatal(err)
		}

		//from install.sh
		// prep data directory
		// mkdir $DATA_DIR/$CCNAME-data
		// mkdir $DATA_DIR/stakewise-data
		// chown $callinguser $DATA_DIR/$CCNAME-data
		// chmod 700 $DATA_DIR/$CCNAME-data
		// # v3-operator user is "nobody" for safety since keys are stored there
		// # you will need to use root to access this directory
		// chown nobody $DATA_DIR/stakewise-data
		u, err := hyperdrive.CallingUser()
		if err != nil {
			log.Errorf("error looking up calling user user info: %e", err)
		}
		os.MkdirAll(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), 0700)
		hyperdrive.Chown(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), u)
		os.MkdirAll(filepath.Join(dataDir, fmt.Sprintf("%s-data", ecName)), 0700)
		hyperdrive.Chown(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), u)
		os.MkdirAll(filepath.Join(dataDir, "stakewise-data"), 0700)
		nobody, err := user.Lookup("nobody")
		hyperdrive.Chown(filepath.Join(dataDir, "stakewise-data"), nobody)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringP("network", "n", "", "Select the network")
	initCmd.Flags().StringP("ecname", "e", "", "Select the execution client [geth, nethermind]")
	initCmd.Flags().String("ccname", "", "Select the consensus client [nimbus]")
	initCmd.Flags().BoolP("remove", "r", false, "Remove the existing installation (if any) in the specified data directory before proceeding with the installation.")
	initCmd.Flags().Bool("checkpoint", false, "Remove the existing installation (if any) in the specified data directory before proceeding with the installation.")

}
