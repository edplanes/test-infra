package cmd

import (
	"log"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete images",
	Long:  `delete images from your hetzner cloud account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := hcloud.NewClient(hcloud.WithToken(token))

		images, err := client.Image.All(cmd.Context())
		if err != nil {
			log.Fatalf("error retrieving images: %v", err)
		}

		if filter != "" {
			images = filterImagesByName(filter, images)
		}

		if olderThan != "" {
			d, err := time.ParseDuration(olderThan)
			if err != nil {
				log.Fatalf("error parsing duration: %v", err)
			}

			images = filterImagesByAge(d, images)
		}

		if userOnly {
			images = filterImagesByUser(images)
		}

		for i, image := range images {
			if i < 2 {
				continue
			}

			log.Printf("deleting image %d", image.ID)
			_, err := client.Image.Delete(cmd.Context(), image)
			if err != nil {
				log.Fatalf("error deleting image %d: %v", image.ID, err)
			}
		}
	},
}
