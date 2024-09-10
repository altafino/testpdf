package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	mindee "github.com/altafino/mindee-client"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Helper function to render error
	renderError := func(message string) {
		tmpl, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct{ ErrorMessage string }{ErrorMessage: message})
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		renderError("Failed to parse form. Please try again.")
		return
	}

	// Get the original filename
	file, header, err := r.FormFile("pdf")
	if err != nil {
		renderError("Failed to get the uploaded file. Please try again.")
		return
	}
	defer file.Close()

	filename := header.Filename
	if filename == "" {
		filename = "Uploaded PDF" // Fallback name if empty
	}

	// Create a temporary file to store the uploaded PDF
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		renderError("Failed to process the uploaded file. Please try again.")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		renderError("Failed to process the uploaded file. Please try again.")
		return
	}

	// Get API key from environment variable
	apiKey := os.Getenv("MINDEE_API_KEY")
	if apiKey == "" {
		renderError("Server configuration error. Please contact the administrator.")
		return
	}

	// Extract data using Mindee API
	invoiceData, err := mindee.GetInvoiceDataForFilePath(tempFile.Name(), apiKey)
	if err != nil {
		renderError("Failed to extract data from the PDF. Please make sure it's a valid invoice and try again.")
		return
	}

	// Marshal the entire invoiceData to JSON
	jsonData, err := json.MarshalIndent(invoiceData, "", "  ")
	if err != nil {
		renderError("Failed to process the extracted data. Please try again.")
		return
	}

	log.Printf("JSON data: %s", string(jsonData)) // Keep this line for logging

	// Create a struct to hold both filename and invoice data
	data := struct {
		Filename    string
		InvoiceData template.JS
	}{
		Filename:    filename,
		InvoiceData: template.JS(jsonData), // Remove the string() conversion
	}

	// Render the result template
	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		renderError("Failed to display the results. Please try again.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		renderError("Failed to display the results. Please try again.")
	}
}
