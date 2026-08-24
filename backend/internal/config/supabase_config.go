package config

import (
	"log"

	storage_go "github.com/supabase-community/storage-go"
)

var SupabaseStorage *storage_go.Client

func InitSupabase(config Config) {
	url := config.Supabase.Url
	key := config.Supabase.Key

	if url == "" || key == "" {
		log.Fatal("SUPABASE_URL or SUPABASE_KEY not set in environment variables")
	}

	storageURL := url + "/storage/v1"

	client := storage_go.NewClient(storageURL, key, nil)
	SupabaseStorage = client

	log.Println("Supabase Storage initialized successfully")
}
