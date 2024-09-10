# PDF Invoice Extractor

This project is a Go HTTP server that hosts an HTMX-based frontend for uploading PDF files and extracting invoice data using the Mindee API with the https://github.com/altafino/mindee-client package.

## Features

- Upload PDF invoices
- Extract data using Mindee API
- Display extracted data in a Tailwind CSS styled HTML page

## Prerequisites

- Go 1.21 or later
- Mindee API key

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/altafino/testpdf.git
   cd testpdf
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Create a `.env` file in the project root and add your Mindee API key:
   ```
   MINDEE_API_KEY=your_api_key_here
   ```

## Running the Application

1. Start the server:
   ```
   go run main.go
   ```

2. Open a web browser and navigate to `http://localhost:8080`

3. Upload a PDF invoice and view the extracted data

## Project Structure

- `main.go`: Contains the main server logic and API endpoints
- `templates/index.html`: The main page HTML template
- `templates/result.html`: The result page HTML template displaying extracted data

## Technologies Used

- Go
- HTMX
- Tailwind CSS
- Mindee API
- altafino/mindee-client
## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
