package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// token is the hetzner cloud api token
	token string

	// olderThan is the age of the image in seconds
	olderThan string

	// filter is a regexp filter for the image name
	filter string

	// userOnly is a flag to only list images created by the user
	userOnly bool
)

var rootCmd = &cobra.Command{
	Use:   "image-cleaner",
	Short: "image-cleaner is a tool to clean up unused images from your hetzner cloud account",
	Long: `image-cleaner is a tool to clean up unused images from your hetzner cloud account.
It is designed to be used as a cronjob.`,
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&token, "token", "t", "", "hetzner cloud api token")
	viper.BindPFlag("HCLOUD_TOKEN", flags.Lookup("token"))
	flags.StringVarP(&olderThan, "older-than", "o", "", "filter images older than")
	flags.StringVarP(&filter, "filter", "f", "", "filter images by name")
	flags.BoolVarP(&userOnly, "user-only", "u", false, "only list user images")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
