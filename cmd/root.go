/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.node-data/nodeset.env)")
	rootCmd.PersistentFlags().StringVarP(&dataDir, "directory", "d", "", "data directory (default is $HOME/.node-data/)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	cfgFile     string
	dataDir     string
	nodesetFile string
)

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if dataDir == "" {
			dataDir = filepath.Join(home, ".node-data")
		}

		if nodesetFile == "" {
			nodesetFile = "nodeset.env"
		}

		viper.AddConfigPath(dataDir)
		viper.SetConfigName(nodesetFile)
		viper.SetConfigType("env")
	}
}
