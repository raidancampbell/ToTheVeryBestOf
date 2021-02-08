package data

import "gorm.io/gorm"

type YoutubeResult struct {
	gorm.Model
	Query   string
	VideoID string
}
