package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type Quote struct {
	Author string
	Quote  string
}

var quotes []Quote

// loadQuotesFromFile loads quotes from a CSV file into the global quotes variable.
func loadQuotesFromFile(filename string) error {
	fmt.Println("Loading quotes from file!")
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Initialize the quotes array
	quotes = make([]Quote, 0, len(records)-1)

	// Iterate over CSV records and create Quote objects
	for i, record := range records {
		// Skip the header row
		if i == 0 {
			continue
		}

		// Create a Quote object and append it to the quotes array
		quote := Quote{
			Author: record[0],
			Quote:  record[1],
		}
		quotes = append(quotes, quote)
	}
	fmt.Printf("Loaded: %d quotes from file.\n", len(records))

	fmt.Println("Loaded quotes from file!")

	return nil
}

func getQuoteByID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	fmt.Printf("Getting quote id: %d\n", id)
	if err != nil || id <= 0 || id > len(quotes)-1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	quote := quotes[id]
	response := map[string]string{
		"author": quote.Author,
		"quote":  quote.Quote,
	}

	fmt.Printf("Returning quote: %v\n", response)
	return c.JSON(response)
}

func getImage(c *fiber.Ctx) error {
	imagePath := "./images/" + "stella.jpg"

	fmt.Println("Getting image", imagePath)

	// Read the image file
	image, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Println("Error reading image:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Set the Content-Type header to inform the client that the response is an image
	c.Set("Content-Type", "image/jpeg") // Adjust the content type based on your image format
	return c.Send(image)
}

func getHello(c *fiber.Ctx) error {
	// Return a simple "Hello World" response
	return c.SendString("Healthy!")
}

func renderHome(c *fiber.Ctx) error {

	// Choose a random quote (you might want to implement a random function)
	quote := quotes[1644]

	fmt.Printf("Quote: %s\n", quote.Quote)
	fmt.Printf("Author: %s\n", quote.Author)

	// Render the HTML template with the quote and author
	return c.Render("index", fiber.Map{
		"Quote":  quote.Quote,
		"Author": quote.Author,
	})
}

func getApp() *fiber.App {
	loadQuotesFromFile("./quotes/quotes.csv")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	// Add the route to the group
	app.Get("/api/v1/image", getImage)
	app.Get("/api/v1/status/health", getHello)
	app.Get("/api/v1/quote/:id", getQuoteByID)

	// Serve static files from the "static" directory
	app.Get("/", renderHome)
	return app
}

func main() {
	app := getApp()
	log.Println("Initialized app")
	app.Listen(":8080")
}
