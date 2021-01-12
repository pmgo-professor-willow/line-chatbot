## Development

```sh
# Listening on port 8080
go run cmd/main.go
```

## Deployment

Prepare project environment variables when GitHub Action do "functions deploy".

### ENV from LINE Console

Visit https://developers.line.biz/console/

* `LINE_CHANNEL_SECRET` & `LINE_CHANNEL_ACCESS_TOKEN`
    * Select a Messaging API channel > Basic settings > Channel secret
    * Select a Messaging API channel > Messaging API > Channel access token

### ENV from GCP

* `GCP_CREDENTIALS`
    * GCP Console > APIs & Services > Credentials > Service Accounts > ADD KEY (JSON)
