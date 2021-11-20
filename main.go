package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const SUBREDDIT = "trains"

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}

	db.AutoMigrate(&submission{})

	// Create a background job that just syncs the latest and posts every hour
	go sync(db)

	// Create HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// need to handle agree/disagree
		}

		sub := random_submission(db)
		log.Println(sub)
	})
	http.ListenAndServe(":8080", nil)
}

func sync(db *gorm.DB) {
	get := func(period string, n int) []*reddit.Post {
		posts, _, err := reddit.DefaultClient().Subreddit.TopPosts(context.Background(), SUBREDDIT, &reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: n,
			},
			Time: period, // hour, day, week, month, year, all
		})

		if err != nil {
			log.Println("", err)
		}

		return posts
	}

	store := func(posts []*reddit.Post) {
		log.Println("have submissions", len(posts))

		for _, post := range posts {
			ext := filepath.Ext(post.URL)
			if ext == "" {
				log.Println("submissions URL isn't direct to file", ext)
				continue
			}

			// this should be noisy as fuck; as most will fail unique constraints
			// and very few will actually be ingested.
			res := db.Create(&submission{
				SubmissionID: post.ID,
				URL:          post.URL,
				Title:        post.Title,
				Author:       post.Author,
			})

			log.Println("stored submission", res.Error, res.RowsAffected)
		}
	}

	if posts := get("all", 10); posts != nil {
		store(posts)
	}

	for {
		// Get the top 25 for the past hour
		if posts := get("hour", 25); posts != nil {
			store(posts)
		}
		<-time.After(1 * time.Hour)
	}
}
