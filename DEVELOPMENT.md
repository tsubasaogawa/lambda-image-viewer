# Development

## Deploying Lambda functions

The Serverless configuration uses `serverless-better-credentials` so that AWS IAM
Identity Center (AWS SSO) profiles work with Serverless Framework v3.

1. Install the viewer's Node.js dependencies:

   ```bash
   cd src/viewer
   npm ci
   ```

2. Authenticate the AWS CLI profile:

   ```bash
   aws sso login --profile default
   ```

3. Deploy with the SSO profile:

   ```bash
   AWS_PROFILE=default AWS_SDK_LOAD_CONFIG=1 sls deploy
   ```

`AWS_SDK_LOAD_CONFIG=1` is required for Serverless to read the SSO profile from
`~/.aws/config`. Run `aws sso login --profile default` again when the SSO session
expires.
