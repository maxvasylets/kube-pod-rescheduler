package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	gitCommit   string
	version     string
	buildDate   string
	buildNumber string
)

var cfgFile string
var showVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-pod-rescheduler",
	Short: "Kubernetes controller that helps to evict and reschedule pods from the node when they're stuck on it by some reasons.",
	Long:  "Kubernetes controller that helps to evict and reschedule pods from the node when they're stuck on it by some reasons.",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("kube-pod-rescheduler\n\nversion: %v\nbuildNumber: %v\ncommit: %v\ndate: %v\n\n", version, buildNumber, gitCommit, buildDate)
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onInitialize)

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version and related information")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.kube-pod-rescheduler.yaml)")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "Don't apply changes to the cluster, just print them")

	// Hide help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

// onInitialize reads in config file and ENV variables if set.
func onInitialize() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kube-pod-rescheduler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kube-pod-rescheduler")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
