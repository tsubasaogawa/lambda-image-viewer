# lambda-image-viewer

A simple image viewer using AWS Lambda + CloudFront + DynamoDB + S3.

## Features

![Screenshot](https://github.com/tsubasaogawa/lambda-image-viewer/assets/7788821/a44fe384-49bc-4b6a-ad44-35d8b331f42a)

- No ALB or API Gateway. The viewer uses Lambda Function URLs.
- The viewer shows EXIF from DynamoDB.
- CloudFront delivery.

![DynamoDB item example](https://github.com/tsubasaogawa/lambda-image-viewer/assets/7788821/c49ad3ba-6123-4196-8cb3-fc9019703229)

## Requirements

- Serverless Framework v3.34
- Terraform v1.5.0
- Go 1.18
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

### Fix environments

```bash
cd lambda-image-viewer

cp -p src/viewer/.env.tmpl src/viewer/.env
vim src/viewer/.env

cp -p terraform/terraform.tfvars.tmpl terraform/terraform.tfvars
vim terraform/terraform.tfvars
```

### Deploy Lambda

```bash
cd src/viewer

make
# Building go binary, deploying lambda using serverless
```

### Deploy instastructure

```bash
cd ../../terraform

terraform init
terraform apply
```

### Upload photo

```bash
aws s3 cp <PHOTO FILE> s3://<CREATED S3 BUCKET>/<FILE NAME>
```

And you should add an item to DynamoDB table created by Terraform.


## Future Works

- Use `provided.al2` Lambda runtime instead of `go1.x`
- Image uploader.
  - Uploading image to S3
  - Adding an item to DynamoDB
