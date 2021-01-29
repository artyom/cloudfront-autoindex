# No more index.html mess with AWS CloudFront/S3

Command cloudfront-autoindex is an AWS Lambda that processes S3 `ObjectCreated`
events looking for objects with `/index.html` suffix in their name, and makes
copies of those objects with `/index.html` and `index.html` suffixes stripped.

This way if you have a directory `doc` with an `index.html` file in it, and you
upload that directory to an S3 bucket fronted by CloudFront, you can then see
your page not only by accessing `https://example.org/doc/index.html`, but also
`https://example.org/doc` and `https://example.org/doc/`, as this lambda will
create two copies of `doc/index.html` key under `doc` and `doc/` keys.

Requires these permissions to work with S3:

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