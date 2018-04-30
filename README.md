# REAME

Tool for uploading cards s3 buckets

## Quickstart

**Install dependencies**

You need the following software to run this application.

```bash
$ go get
```

Create a `.env` file in the root directory and set the following env variables in the file.

```
BUCKET s3 bucket name to upload cards to
BASE_PATH Path where cards repository is located
USE_LOCAL Set to true if you are want to copy cards to a local dir (for local development)
TARGET_PATH Path to copy cards to (for local develoment)

```
