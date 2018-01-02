#!/bin/sh

# just some settings I use for local development

export LUMEN_PORT=9989
export LUMEN_SECRET=notreallysecret
export LUMEN_S3_BUCKET=thraxil-lumen-dev

go build .
./lumen
