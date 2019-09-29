<h1 align="center">s3-proxy</h1>

<p align="center">
A proxying server to private buckets in S3
</p>

### Introduction

There are many use cases where S3 is used as an object store for objects that may be intended to be accessed publicly.
Sometimes it is a requirement that restrictions be placed on who can access those objects without using the S3 API (eg. an company internal static site).
Since AWS does not provide the tools to do this, s3-proxy was born.

s3-proxy is meant to be completely configuration driven so that no source code modification or forking is necessary.
It can be deployed to your own private servers or a platform like Heroku with ease.
It supports basic auth for the use case of deploying to a publicly accessible server, although it is recommended to deploy s3-proxy within a firewall.

### Configuration

s3-proxy can be configured with yaml file. You can place your `config.yml` file in `/etc/s3-proxy`, `$HOME/.s3-proxy/` or same folder as binary.
For configuration of aws s3 credentials you can look official [documentation](https://docs.aws.amazon.com/en_us/sdk-for-go/v1/developer-guide/configuring-sdk.html).

Sample configuration:

```yaml
sites:
  - host: foo.com
    bucket: foo
    options:
      website: true
      gzip: false
      cors: true
```

### A note about AWS keys

It is good practice to utilize proper user management with the keys that are deployed with s3-proxy.
Any keys are that are used for proxying should be limited to have read-only access to the S3 buckets that they intend to fetch from.
Read-only access translates to the permissions: s3:GetObject, s3:GetBucketWebsite, s3:ListBucket.
