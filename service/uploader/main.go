package main

import (
	"e-commerce-microservice/uploader/storage"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	supabaseURL := os.Getenv("SUPABASE_URL") // Ganti dengan .env atau hardcoded jika mau
	supabaseKey := os.Getenv("SUPABASE_KEY")
	bucket := "your-bucket-name" // Ganti dengan nama bucket kamu

	uploader, err := storage.NewSupabaseUploader(supabaseURL, supabaseKey, bucket)
	if err != nil {
		log.Fatalf("❌ Failed to create Supabase client: %v", err)
	}

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		publicURL, err := uploader.UploadImage(r.Context(), header, "car")
		if err != nil {
			http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim response JSON manual
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"url":"%s"}`, publicURL)
	})
	fmt.Println("✅ Uploader service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
