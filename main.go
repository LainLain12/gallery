package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Wallpaper struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	ImageURL   string   `json:"imageUrl"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	Resolution string   `json:"resolution"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    []Wallpaper `json:"data"`
	Message string      `json:"message,omitempty"`
}

var (
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".webp", ".bmp"}
	categories      = []string{"nature", "culture", "digital"}
	resolutions     = []string{"1080x1920", "1440x2560", "2160x3840", "1080x2340", "1170x2532"}
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create Gin router
	r := gin.Default()

	// Enable CORS for Android app
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/wallpapers/:category", getWallpapersByCategory)
		api.GET("/wallpapers/:category/random", getRandomWallpaper)
		api.GET("/wallpapers", getAllWallpapers)
		api.GET("/categories", getCategories)
	}

	// Serve static images
	r.Static("/images", "./images")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Wallpaper API is running"})
	})

	// Start server
	fmt.Println("🚀 Wallpaper API Server starting on http://localhost:8664")
	fmt.Println("📁 Place your images in:")
	fmt.Println("   - images/nature/")
	fmt.Println("   - images/culture/")
	fmt.Println("   - images/digital/")
	fmt.Println("📡 API Endpoints:")
	fmt.Println("   - GET /api/v1/wallpapers/nature")
	fmt.Println("   - GET /api/v1/wallpapers/culture")
	fmt.Println("   - GET /api/v1/wallpapers/digital")
	fmt.Println("   - GET /api/v1/wallpapers/{category}/random")

	log.Fatal(r.Run(":8664"))
}

func getWallpapersByCategory(c *gin.Context) {
	category := c.Param("category")

	// Validate category
	if !isValidCategory(category) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid category. Use: nature, culture, or digital",
		})
		return
	}

	wallpapers, err := loadWallpapersFromFolder(c, category)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Error loading wallpapers: %v", err),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    wallpapers,
	})
}

func getRandomWallpaper(c *gin.Context) {
	category := c.Param("category")

	// Validate category
	if !isValidCategory(category) {
		c.JSON(400, APIResponse{
			Success: false,
			Message: "Invalid category. Use: nature, culture, or digital",
		})
		return
	}

	wallpapers, err := loadWallpapersFromFolder(c, category)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Error loading wallpapers: %v", err),
		})
		return
	}

	if len(wallpapers) == 0 {
		c.JSON(404, APIResponse{
			Success: false,
			Message: "No wallpapers found in this category",
		})
		return
	}

	// Get random wallpaper
	randomIndex := rand.Intn(len(wallpapers))
	randomWallpaper := wallpapers[randomIndex]

	c.JSON(200, APIResponse{
		Success: true,
		Data:    []Wallpaper{randomWallpaper},
	})
}

func getAllWallpapers(c *gin.Context) {
	allWallpapers := []Wallpaper{}

	for _, category := range categories {
		wallpapers, err := loadWallpapersFromFolder(c, category)
		if err != nil {
			log.Printf("Error loading %s wallpapers: %v", category, err)
			continue
		}
		allWallpapers = append(allWallpapers, wallpapers...)
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    allWallpapers,
	})
}

func getCategories(c *gin.Context) {
	c.JSON(200, gin.H{
		"success":    true,
		"categories": categories,
	})
}

func loadWallpapersFromFolder(c *gin.Context, category string) ([]Wallpaper, error) {
	folderPath := filepath.Join("images", category)

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %v", folderPath, err)
	}

	var wallpapers []Wallpaper
	id := 1

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Check if file is an image
		if !isImageFile(file.Name()) {
			continue
		}

		// Generate dynamic base URL from request
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

		// Generate wallpaper data
		wallpaper := Wallpaper{
			ID:         id,
			Title:      generateRandomTitle(category),
			ImageURL:   fmt.Sprintf("%s/images/%s/%s", baseURL, category, file.Name()),
			Category:   strings.Title(category),
			Tags:       generateRandomTags(category),
			Resolution: getRandomResolution(),
		}

		wallpapers = append(wallpapers, wallpaper)
		id++
	}

	return wallpapers, nil
}

func isValidCategory(category string) bool {
	for _, validCategory := range categories {
		if category == validCategory {
			return true
		}
	}
	return false
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range imageExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func generateRandomTitle(category string) string {
	titles := map[string][]string{
		"nature": {
			"Serene Landscape", "Mountain Vista", "Ocean Breeze", "Forest Path",
			"Sunset Glory", "River Flow", "Desert Bloom", "Alpine View",
			"Coastal Beauty", "Wilderness", "Garden Paradise", "Peaceful Lake",
		},
		"culture": {
			"Ancient Heritage", "Traditional Art", "Cultural Festival", "Historic Monument",
			"Ethnic Pattern", "Sacred Temple", "Folk Design", "Heritage Site",
			"Cultural Symbol", "Traditional Craft", "Ancient Wisdom", "Cultural Legacy",
		},
		"digital": {
			"Cyber Grid", "Digital Wave", "Neon Dreams", "Tech Pattern",
			"Futuristic Design", "Digital Art", "Cyber Space", "Modern Abstract",
			"Tech Innovation", "Digital Future", "Cyber Aesthetic", "Virtual Reality",
		},
	}

	categoryTitles := titles[category]
	if len(categoryTitles) == 0 {
		return "Beautiful Wallpaper"
	}

	return categoryTitles[rand.Intn(len(categoryTitles))]
}

func generateRandomTags(category string) []string {
	tagSets := map[string][]string{
		"nature":  {"landscape", "natural", "scenic", "outdoor", "peaceful", "green", "blue", "mountains", "ocean", "forest"},
		"culture": {"traditional", "heritage", "ancient", "artistic", "cultural", "historic", "ethnic", "sacred", "folk", "classic"},
		"digital": {"modern", "futuristic", "tech", "cyber", "digital", "abstract", "neon", "geometric", "virtual", "electronic"},
	}

	availableTags := tagSets[category]
	if len(availableTags) == 0 {
		return []string{"wallpaper"}
	}

	// Select 2-4 random tags
	numTags := rand.Intn(3) + 2 // 2 to 4 tags
	if numTags > len(availableTags) {
		numTags = len(availableTags)
	}

	// Shuffle and select tags
	shuffled := make([]string, len(availableTags))
	copy(shuffled, availableTags)

	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:numTags]
}

func getRandomResolution() string {
	return resolutions[rand.Intn(len(resolutions))]
}
