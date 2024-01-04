/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nodeset-org/hyperdrive-stakewise/config"
	"github.com/nodeset-org/hyperdrive-stakewise/local"
)

// var Local embed.FS

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalizes the ~/.node-data/ directory with nodeset.env, compose.yaml and the ec and cc docker files.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("init called")
		network, _ := cmd.Flags().GetString("network")
		ecName, _ := cmd.Flags().GetString("ecname")
		if ecName == "" {
			ecs := []string{"geth", "nethermind"}
			n := rand.Int() % len(ecs)
			ecName = ecs[n]
		}

		ccName, _ := cmd.Flags().GetString("ccname")
		if ccName == "" {
			ccs := []string{"nimbus"}
			n := rand.Int() % len(ccs)
			ccName = ccs[n]
		}

		// var envFile []byte
		var err error
		var c config.Config
		switch network {
		case "holskey":
			c = config.Holskey
		case "holskey-dev":
			c = config.HoleskyDev
		case "main":
			c = config.Gravita

		default:
			log.Fatalf("network %s is not avaliable, please choose holskey, holskey-dev or main", network)
		}

		err = os.MkdirAll(dataDir, 0766)
		if err != nil {
			log.Error(err)
		}

		c.SetViper()
		//Ensure that nodeset.env contains the correct ECNAME and CCNAME
		viper.Set("ECNAME", ecName)
		viper.Set("CCNAME", ccName)

		err = c.SetConfigPath(dataDir)
		if err != nil {
			log.Fatal(err)
		}
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

		// //from install.sh
		// // prep data directory
		// // mkdir $DATA_DIR/$CCNAME-data
		// // mkdir $DATA_DIR/stakewise-data
		// // chown $callinguser $DATA_DIR/$CCNAME-data
		// // chmod 700 $DATA_DIR/$CCNAME-data
		// // # v3-operator user is "nobody" for safety since keys are stored there
		// // # you will need to use root to access this directory
		// // chown nobody $DATA_DIR/stakewise-data
		// os.MkdirAll(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), 0766)
		// // u, err := user.Current()
		// if err != nil {
		// 	log.Errorf("error looking up current user user info: %e", err)
		// }
		// // chown(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), u)

		// os.MkdirAll(filepath.Join(dataDir, fmt.Sprintf("%s-data", ecName)), 0766)
		// // chown(filepath.Join(dataDir, fmt.Sprintf("%s-data", ccName)), u)
		// os.MkdirAll(filepath.Join(dataDir, "stakewise-data"), 0766)

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
	initCmd.Flags().StringP("network", "n", "holskey", "Select the network [main, holskey is Default]")
	initCmd.Flags().StringP("ecname", "e", "", "Select the execution client [geth, nethermind]")
	initCmd.Flags().StringP("ccname", "c", "", "Select the consensus client [nimbus]")

}

func chown(dir string, u *user.User) error {

	if runtime.GOOS != "windows" {
		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)

		err := syscall.Chown(dir, uid, gid)
		if err != nil {
			return err
		}
	}
	return nil
}
