Here are important commands for developing, testing, formatting, and deploying code in this project:

### General System Utilities

*   `git`: Version control system.
*   `ls`: List directory contents.
*   `cd`: Change directory.
*   `grep`: Search for patterns in files.
*   `find`: Search for files in a directory hierarchy.

### Version Management (Recommended)

*   `slsenv`: Manage Serverless Framework versions.
*   `tfenv`: Manage Terraform versions.
*   `goenv`: Manage Go versions.

### Go Specific Commands

*   **Testing:** `go test ./...` (runs all tests in the current Go module)
*   **Formatting:** `go fmt ./...` (formats all Go files in the current module)

### Deployment Commands

**1. Deploy Lambda Functions:**

```bash
cd src/viewer
cp -p .env.tmpl .env
vim .env # Configure environment variables
serverless deploy
serverless info # Show deployment information
```

**2. Deploy AWS Infrastructure (Terraform):**

```bash
cd terraform
cp -p terraform.tfvars.tmpl terraform.tfvars
vim terraform.tfvars # Configure Terraform variables
terraform init
terraform apply
```

**3. Re-deploy Lambda after Infrastructure changes (if needed):**

```bash
cd src/viewer
serverless deploy
```

**4. Upload Photo to S3:**

```bash
aws s3 cp <PHOTO FILE> s3://<CREATED S3 BUCKET>/<FILE NAME>
```