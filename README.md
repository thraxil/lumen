# lumen

lightweight image serving

Basically, a trimmed down version of
[reticulum](https://github.com/thraxil/reticulum/) without the
distributed aspect. Sometimes it's enough to essentially have a single
node backed by S3.

The public API (upload images and retrieve scaled images by URL) is
fully compatible with reticulum (so existing client code can just
switch). But all images are stored in a single S3 bucket. For low
amounts of traffic, as long as you are happy with S3's durability
guarantees, this could be perfectly fine for you, and it's simpler to
setup and run than a full reticulum cluster (which was originally
designed and optimzed for running on commodity hardware and drives).
