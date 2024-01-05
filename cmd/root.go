/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-stakewise/hyperdrive/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	// Used for flags.
	dataDir string
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "hyperdrive-stakewise",
		Short: "A brief description of your application",
		Long: `{::} NodeSet Hyperdrive - StakeWise | VERSION {::}

	Usage: 
	
		nodeset [OPTIONS] COMMAND
	
	Options: 
	
		-h, --help
	
			Show this message
	
		-d directory, --data-directory=directory
	
			Specify location for the configuration directory. Default is /home/$USER/.node-data.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	//if user is sudo, use the calling user home
	var callingUser string
	if os.Geteuid() == 0 { //sudo
		callingUser = os.Getenv("SUDO_USER")
		user, err := user.Lookup(callingUser)
		if err != nil {
			log.Fatal(err)
		}
		dirname = user.HomeDir
	}

	d := filepath.Join(dirname, ".node-data")
	rootCmd.PersistentFlags().StringVarP(&dataDir, "directory", "d", d, "data directory")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "nodeset.env", "config file name")
	config.ConfigFile = cfgFile
	//Set the global DataDir and config path based on the gloabal flag
	config.SetConfigPath(dataDir)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	config.ConfigFile = cfgFile
	config.SetConfigPath(dataDir)
}
