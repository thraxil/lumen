#!/bin/sh

# just some settings I use for local development

export LUMEN_PORT=9989
export LUMEN_SECRET=notreallysecret
export LUMEN_S3_BUCKET=thraxil-lumen-dev
export LUMEN_S3_ACCESS_KEY=$AWS_ACCESS_KEY
export LUMEN_S3_SECRET_KEY=$AWS_SECRET_KEY

go build .
./lumen
