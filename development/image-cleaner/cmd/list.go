package cmd

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var (
	// olderThan is the age of the image in seconds
	olderThan string

	// filter is a regexp filter for the image name
	filter string

	// userOnly is a flag to only list images created by the user
	userOnly bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list images",
	Long:  `list images from your hetzner cloud account.`,
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

		for _, image := range images {
			fmt.Printf("%d|%s|%s|%s\n", image.ID, image.Name, image.Type, image.Status)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&olderThan, "older-than", "o", "", "filter images older than")
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "filter images by name")
	listCmd.Flags().BoolVarP(&userOnly, "user-only", "u", false, "only list user images")
}

func filterImagesByName(filter string, images []*hcloud.Image) []*hcloud.Image {
	var filteredImages []*hcloud.Image

	for _, image := range images {
		re := regexp.MustCompile(filter)
		if re.MatchString(image.Name) {
			filteredImages = append(filteredImages, image)
		}
	}

	return filteredImages
}

func filterImagesByAge(olderThan time.Duration, images []*hcloud.Image) []*hcloud.Image {
	if olderThan == 0 {
		return images
	}

	var filteredImages []*hcloud.Image

	for _, image := range images {
		if image.Created.Before(time.Now().Add(-olderThan)) {
			filteredImages = append(filteredImages, image)
		}
	}

	return filteredImages
}

func filterImagesByUser(images []*hcloud.Image) []*hcloud.Image {
	var filteredImages []*hcloud.Image

	for _, image := range images {
		if image.Type == hcloud.ImageTypeSnapshot {
			filteredImages = append(filteredImages, image)
		}
	}

	return filteredImages
}
