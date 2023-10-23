# lambda-image-viewer

A simple image viewer using AWS Lambda + CloudFront + DynamoDB + S3.

## Features

![Screenshot](https://github.com/tsubasaogawa/lambda-image-viewer/assets/7788821/ec35bdf9-1446-4f5c-a85e-3a82940aeef3)

- No ALB or API Gateway. The viewer uses Lambda Function URLs.
- The viewer shows some EXIFs from DynamoDB.
- S3-event-driven Tagger gets EXIF by a photo and puts to DynamoDB.
- CloudFront delivery.

![DynamoDB item example](https://github.com/tsubasaogawa/lambda-image-viewer/assets/7788821/3ff31067-5d92-4d71-8bb5-8b2568558fc8)


## Requirements

- Serverless Framework v3.34
- Terraform v1.5.0
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

## Future Works

- Fix design
