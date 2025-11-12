The project has the following rough structure:

*   `src/viewer`: Contains the Go Lambda functions and Serverless configuration.
    *   `cmd/cleanup`: Lambda function for cleanup tasks.
    *   `cmd/tagger`: Lambda function for extracting EXIF data from images.
    *   `cmd/thumbnail`: Lambda function for generating image thumbnails.
    *   `cmd/viewer`: Main Lambda function for the image viewer.
    *   `internal/model`: Go structs defining data models (e.g., `item.go`, `metadata.go`).
    *   `templates`: HTML templates for the image viewer (e.g., `camera_roll.html.tmpl`, `index.html.tmpl`).
*   `terraform`: Contains Terraform configurations for provisioning AWS infrastructure.
    *   `assets`: JavaScript files used by Terraform for Lambda@Edge functions or other purposes.
*   `docs`: Contains project documentation and diagrams (e.g., `diagram.drawio.png`).
*   `.gemini/`, `.serena/`: Configuration and memory files for the Gemini CLI agent.
*   `.clinerules`, `.gitignore`, `LICENSE`, `README.md`: Project-level configuration and documentation files.