package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute: %s %v", command, args)
	}
}

func createFile(path string, content string) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		log.Fatalf("Cannot create directory %s: %v", path, err)
	}
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Cannot create file %s: %v", path, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Cannot write to file %s: %v", path, err)
	}
}

func setupProject(projectName string) {
	fmt.Printf("Setting up your project: %s...\n", projectName)

	// Create project directory and initialize module
	runCommand("mkdir", projectName)

	// Change to project directory
	err := os.Chdir(projectName)
	if err != nil {
		log.Fatalf("Failed to change directory to %s: %v", projectName, err)
	}

	runCommand("go", "mod", "init", projectName)

	// Install dependencies
	runCommand("go", "get", "-u", "github.com/gin-gonic/gin")
	runCommand("go", "get", "-u", "gorm.io/gorm")
	runCommand("go", "get", "-u", "gorm.io/driver/postgres")
	runCommand("go", "get", "-u", "github.com/gin-contrib/cors")
	runCommand("go", "get", "-u", "github.com/joho/godotenv")

	// Create .env file
	envContent := fmt.Sprintf("DATABASE_URL=host=localhost user=postgres dbname=%s_db sslmode=disable password=yourpassword", projectName)
	createFile(".env", envContent)

	// Create main.go
	mainContent := `package main

import (
    "net/http"
    "` + projectName + `/models"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Initialize Gin router
    router := gin.Default()

    // Setup CORS
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }))

    // Initialize GORM with PostgreSQL database
    dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Auto migrate models
    err = db.AutoMigrate(&models.User{})
    if err != nil {
        panic("failed to migrate database")
    }

    // Define a simple route
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to the Gin-GORM project!",
        })
    })

    // Start the server
    router.Run(":8080")
}
`
	createFile("main.go", mainContent)

	// Create models directory and user.go
	modelContent := `package models

import (
    "gorm.io/gorm"
)

// User represents a user model
type User struct {
    gorm.Model
    Name  string ` + "`json:\"name\"`" + `
    Email string ` + "`json:\"email\" gorm:\"unique\"`" + `
}
`
	createFile("models/user.go", modelContent)

	fmt.Println("Project setup completed successfully.")
}

func main() {
	app := &cli.App{
		Name:  "go-create-project",
		Usage: "Create a new Go project with Gin and GORM",
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "create a new project",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("provide a project name")
					}
					projectName := c.Args().Get(0)
					setupProject(projectName)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
