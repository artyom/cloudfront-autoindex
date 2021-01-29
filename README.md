# No more index.html mess with AWS CloudFront/S3

Command cloudfront-autoindex is an AWS Lambda that processes S3 `ObjectCreated`
events looking for objects with `/index.html` suffix in their name, and makes
copies of those objects with `/index.html` and `index.html` suffixes stripped.

This way if you have a directory `doc` with an `index.html` file in it, and you
upload that directory to an S3 bucket fronted by CloudFront, you can then see
your page not only by accessing `https://example.org/doc/index.html`, but also
`https://example.org/doc` and `https://example.org/doc/`, as this lambda will
create two copies of `doc/index.html` key under `doc` and `doc/` keys.

## Setup

Build and compress:

    GOOS=linux GOARCH=amd64 go build -o main
    zip -9 lambda.zip main

Create a new AWS Lambda, picking the "Go 1.x" runtime. Change its handler name from default "hello" to "main" (binary name you built above), and upload lambda.zip file.

It requires the usual permissions, e.g. `AWSLambdaBasicExecutionRole` AWS-managed role, plus these permissions to work with S3 (optionally limit with your bucket names):

    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "CopyObject",
                "Effect": "Allow",
                "Action": [
                    "s3:PutObject",
                    "s3:GetObject"
                ],
                "Resource": "arn:aws:s3:::*/*"
            }
        ]
    }

Once Lambda is created and configured, go to S3 bucket settings, Properties → Event notifications → Create event notification. Enter `index.html` as a *Suffix*, and select `s3:ObjectCreated` events except for the `s3:ObjectCreated:Copy`. Pick your Lambda function created above as event destination.
