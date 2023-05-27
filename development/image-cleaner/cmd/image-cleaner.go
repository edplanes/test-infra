package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// token is the hetzner cloud api token
	token string
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

	rootCmd.AddCommand(listCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
