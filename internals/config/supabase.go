package config

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

type SupabaseConfig struct {
	URL    string
	Key    string
	Client *supabase.Client
}

func GetSupabaseURL() string {
	fmt.Println("Getting Supabase URL...")
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		fmt.Println("[ERROR] SUPABASE_URL not found!")
	} else {
		fmt.Println("[INFO] SUPABASE_URL successfully loaded:", maskURL(supabaseURL))
	}

	return supabaseURL
}

func GetSupabaseKey() string {
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseKey == "" {
		fmt.Println("[ERROR] SUPABASE_KEY not found!")
	} else {
		fmt.Println("[INFO] SUPABASE_KEY successfully loaded:", maskKey(supabaseKey))
	}
	return supabaseKey
}

func maskURL(url string) string {
	if len(url) <= 12 {
		return "***"
	}
	visible := url[:8]
	return visible + "***"
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	visible := key[:4]
	return visible + "***"
}

func NewSupabaseClient(supabaseURL, supabaseKey string) (*SupabaseConfig, error) {
	fmt.Println("[INFO] Mencoba menghubungkan ke Supabase...")
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		fmt.Println("[ERROR] Failed to connect to Supabase:", err.Error())
		return nil, err
	}
	fmt.Println("[SUCCESS] Successfully connected to Supabase")

	return &SupabaseConfig{
		URL:    supabaseURL,
		Key:    supabaseKey,
		Client: client,
	}, nil
}

func InitSupabaseClient() (*SupabaseConfig, error) {
	fmt.Println("[INFO] Starting Supabase initialization...")
	config, err := NewSupabaseClient(GetSupabaseURL(), GetSupabaseKey())
	if err != nil {
		return nil, err
	}
	fmt.Println("[INFO] Supabase initialization completed")
	return config, nil
}

func (s *SupabaseConfig) GetClient() *supabase.Client {
	return s.Client
}
