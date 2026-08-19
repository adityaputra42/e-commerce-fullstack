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
		log.Println("⚠️  Supabase Storage not configured (SUPABASE_URL/SUPABASE_KEY missing) — " +
			"image upload endpoints will return an error until these are set. Everything else will run normally.")
		return
	}

	storageURL := url + "/storage/v1"

	client := storage_go.NewClient(storageURL, key, nil)
	SupabaseStorage = client

	log.Println("Supabase Storage initialized successfully")
}
