package main

import (
	"e-commerce-microservice/uploader/config"
	"e-commerce-microservice/uploader/storage"
	"fmt"

	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
		panic(err)
	}

	uploader, err := storage.NewSupabaseUploader(conf.SupabaseUrl, conf.SupabaseKey, conf.Bucket)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Supabase client")

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

		form := r.MultipartForm
		files := form.File["images"]
		if len(files) == 0 {
			// Fall back to single file upload if no multiple files found
			file, header, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "No files uploaded", http.StatusBadRequest)
				return
			}
			defer file.Close()

			publicURL, err := uploader.UploadImage(r.Context(), header, "car")
			if err != nil {
				http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"urls":["%s"]}`, publicURL)
			return
		}

		urls, err := uploader.UploadImages(r.Context(), files, "car")
		if err != nil {
			http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert urls slice to JSON array string
		urlJSON := "["
		for i, url := range urls {
			if i > 0 {
				urlJSON += ","
			}
			urlJSON += fmt.Sprintf(`"%s"`, url)
		}
		urlJSON += "]"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"urls":%s}`, urlJSON)
	})
	fmt.Println("✅ Uploader service started on :8080")
	http.ListenAndServe(":8080", nil)
	log.Info().Msg("Uploader service started on :8080")
}
