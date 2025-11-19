# lambda-image-viewer

A simple image viewer using AWS Lambda + CloudFront + DynamoDB + S3.

I used to use Flickr as a storage for blog images, but for several reasons (mainly financial), I needed an alternative. Therefore, I am prioritizing the implementation of features that are essential for my personal use.

## Features

### CDN

route: `<CDN Domain>/.../foo.jpg`

- Provides a standard S3-backed CDN via CloudFront.
- Private mode: Images in the `private` folder are accessible only with the appropriate token query.

### Image Viewer

route: `<Viewer Domain>/image/.../foo.jpg`

<img width="2280" height="1512" alt="image" src="https://github.com/user-attachments/assets/b28f7289-b46f-4cb3-a73a-aeeef89f45bb" />

- Shows an image info (EXIF).

### Camera Roll (Admin Page)

route: `<Viewer Domain>/cameraroll/`

<img width="1299" height="992" alt="image" src="https://github.com/user-attachments/assets/eb065214-bbe1-4a2e-8f6a-e5ef5152291b" />

- You can list all uploaded images.
- Basic auth.
- Private images (`/private/*.jpg`): only accessible from camera roll page.

## Structure

![Diagram](./docs/diagram.drawio.png)

- The viewer does not use ALB but Lambda Function URLs.
- S3-event-driven Tagger gets EXIF by a photo and puts to DynamoDB.

![DynamoDB item example](https://github.com/tsubasaogawa/lambda-image-viewer/assets/7788821/3ff31067-5d92-4d71-8bb5-8b2568558fc8)

## Requirements

- Serverless Framework v3.34
- Terraform v1.12 or above
- Go 1.21
- AWS

You can also use **env to install easily.

- [tsubasaogawa/slsenv](https://github.com/tsubasaogawa/slsenv)
- [tfutils/tfenv](https://github.com/tfutils/tfenv)
- [go-nv/goenv](https://github.com/go-nv/goenv)

## Usage

### Clone

```bash
git clone https://github.com/tsubasaogawa/lambda-image-viewer.git
```

### Deploy Lambda

```bash
cd src/viewer

# Fix environments. Please follow comments in .env file.
cp -p .env.tmpl .env
vim .env

# Building go binary, deploying lambda
serverless deploy

# Show the deploy information.
serverless info
```

### Deploy Infrastructure

```bash
cd ../../terraform

# Fix environments.
cp -p terraform.tfvars.tmpl terraform.tfvars
vim terraform.tfvars

terraform init
terraform apply

# Deploy Lambda again to set some environment variables using Secrets Manager created by Terraform.
cd ../src/viewer
serverless deploy
```

### Upload photo

```bash
aws s3 cp <PHOTO FILE> s3://<CREATED S3 BUCKET>/<FILE NAME>
```

## Future Works

- Fix design
