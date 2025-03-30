package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary

func InitCloudinary() {
	var err error
	Cloudinary, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("❌ Failed to initialize Cloudinary: %v", err)
	}
	log.Println("✅ Cloudinary initialized successfully!")
}
