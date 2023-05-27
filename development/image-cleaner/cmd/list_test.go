package cmd

import (
	"reflect"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func Test_filterImagesByName(t *testing.T) {
	tests := []struct {
		name   string
		filter string
		images []*hcloud.Image
		want   []*hcloud.Image
	}{
		{
			name:   "empty filter",
			filter: "",
			images: []*hcloud.Image{
				{
					ID:   1,
					Name: "test",
				},
				{
					ID:   2,
					Name: "test2",
				},
			},
			want: []*hcloud.Image{
				{
					ID:   1,
					Name: "test",
				},
				{
					ID:   2,
					Name: "test2",
				},
			},
		},
		{
			name:   "filter",
			filter: "test",
			images: []*hcloud.Image{
				{
					ID:   1,
					Name: "test",
				},
				{
					ID:   2,
					Name: "test2",
				},
				{
					ID:   3,
					Name: "tes3",
				},
			},
			want: []*hcloud.Image{
				{
					ID:   1,
					Name: "test",
				},
				{
					ID:   2,
					Name: "test2",
				},
			},
		},
		{
			name:   "complex filter",
			filter: "db-ubuntu-22",
			images: []*hcloud.Image{
				{
					ID:   1,
					Name: "db-ubuntu-20.04",
				},
				{
					ID:   2,
					Name: "db-ubuntu-22.04",
				},
				{
					ID:   3,
					Name: "k8s-ubuntu-20.04",
				},
			},
			want: []*hcloud.Image{
				{
					ID:   2,
					Name: "db-ubuntu-22.04",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterImagesByName(tt.filter, tt.images); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterImagesByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterImagesByAge(t *testing.T) {
	now := time.Now()
	hourOld := now.Add(-time.Hour)
	thrityMinutesOld := now.Add(-time.Minute * 30)
	dayOld := now.Add(-time.Hour * 24)

	tests := []struct {
		name      string
		olderThan time.Duration
		images    []*hcloud.Image
		want      []*hcloud.Image
	}{
		{
			name:      "return all images",
			olderThan: time.Duration(0),
			images: []*hcloud.Image{
				{
					ID:      1,
					Name:    "test",
					Created: now,
				},
				{
					ID:      2,
					Name:    "test2",
					Created: hourOld,
				},
			},
			want: []*hcloud.Image{
				{
					ID:      1,
					Name:    "test",
					Created: now,
				},
				{
					ID:      2,
					Name:    "test2",
					Created: hourOld,
				},
			},
		},
		{
			name:      "return images older than 1 hour",
			olderThan: time.Hour,
			images: []*hcloud.Image{
				{
					ID:      1,
					Name:    "test",
					Created: now,
				},
				{
					ID:      2,
					Name:    "test2",
					Created: thrityMinutesOld,
				},
				{
					ID:      3,
					Name:    "test3",
					Created: dayOld,
				},
			},
			want: []*hcloud.Image{
				{
					ID:      3,
					Name:    "test3",
					Created: dayOld,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterImagesByAge(tt.olderThan, tt.images); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterImagesByAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterImagesByUser(t *testing.T) {
	tests := []struct {
		name   string
		images []*hcloud.Image
		want   []*hcloud.Image
	}{
		{
			name: "only snapshots",
			images: []*hcloud.Image{
				{
					ID:   1,
					Type: hcloud.ImageTypeSnapshot,
				},
				{
					ID:   2,
					Type: hcloud.ImageTypeSnapshot,
				},
				{
					ID:   3,
					Type: hcloud.ImageTypeSnapshot,
				},
			},
			want: []*hcloud.Image{
				{
					ID:   1,
					Type: hcloud.ImageTypeSnapshot,
				},
				{
					ID:   2,
					Type: hcloud.ImageTypeSnapshot,
				},
				{
					ID:   3,
					Type: hcloud.ImageTypeSnapshot,
				},
			},
		},
		{
			name: "return only snapshots",
			images: []*hcloud.Image{
				{
					ID:   1,
					Type: hcloud.ImageTypeSnapshot,
				},
				{
					ID:   2,
					Type: hcloud.ImageTypeSystem,
				},
				{
					ID:   3,
					Type: hcloud.ImageTypeApp,
				},
			},
			want: []*hcloud.Image{
				{
					ID:   1,
					Type: hcloud.ImageTypeSnapshot,
				},
			},
		},
		{
			name: "no snapshots",
			images: []*hcloud.Image{
				{
					ID:   2,
					Type: hcloud.ImageTypeSystem,
				},
				{
					ID:   3,
					Type: hcloud.ImageTypeApp,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterImagesByUser(tt.images); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterImagesByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
