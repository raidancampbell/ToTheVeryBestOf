package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raidancampbell/ToTheVeryBestOf/data"
	"github.com/raidancampbell/ToTheVeryBestOf/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

// requirements:
// YOUTUBE_API_KEY env variable set to a youtube data v3 api key
// LASTFM_API_KEY env variable set to a last.fm API key
// defined up here to clarify the requirements
var (
	youtubeKey = os.Getenv("YOUTUBE_API_KEY")
	lastfmKey  = os.Getenv("LASTFM_API_KEY")
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("best.db"),nil)
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	// Migrate the schema
	db.AutoMigrate(&data.YoutubeResult{})
}

func main() {
	topTracksHandler := handlers.NewTopTracks(youtubeKey, lastfmKey, db)
	r := gin.Default()
	r.GET("/", handlers.Landing)
	r.GET("/artist", topTracksHandler.HandleArtistRequest)
	_ = r.Run("0.0.0.0:9071")
}
