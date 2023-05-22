# paczkobot

A Telegram bot for tracking packages.

See it in action: [@paczko_bot](https://t.me/paczko_bot)

## Features
- [x] Track packages from these providers
  - [x] dhl (requires API key)
  - [x] ups
  - [x] dpd.com.pl
  - [x] poczta-polska
  - [x] postnl
  - [x] Inpost
  - [x] gls
  - [x] packeta
- Follow packages and send notifications when a package status changes
- Generate QR codes for InPost Paczkomaty
- Automatically import packages from InPost Paczkomaty
- Remotely open InPost Paczkomaty
- Detect package barcodes in images


## Screenshots

![Screenshot](./docs/tracking.jpg)
![Screenshot](./docs/barcode.png)
![Screenshot](./docs/inpostopen.png)
## Usage

Create a file called `paczkobot.yaml` looking like this:

```yaml
telegram:
  debug: false
  username: paczko_bot
  token: "<telegram api key>"
db:
  type: sqlite # or postgres
  filename: paczkobot_dev.db # for sqlite
  dsn: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Europe/Warsaw" # for postgres
tracking:
  providers:
    dhl:
      enable: false # register on the dhl developers webpage to enable
      api_key: "<dhl api key>"
    mock:
      enable: false
  automatic_tracking_check_interval: 20m0s
  automatic_tracking_check_jitter: 7m0s
  delay_between_packages_in_automatic_tracking: 1m0s
  max_packages_per_automatic_tracking_check: 15
  max_time_without_change: 336h0m0s

```

You have to enter your telegram token there.

Then you can run the bot as a Go program (`go run .`) or use the following `docker-compose.yml` file to run in docker:

```
version: '3'
services:
  paczkobot:
    image: ghcr.io/alufers/paczkobot:latest
    restart: unless-stopped
    volumes:
        - ./paczkobot-config.yaml:/etc/paczkobot/paczkobot.yaml
```

Images for x86_64 and aarch64 are provided.

## Contributing

Contributions are welcome! Please open an issue or a pull request.

Please note that this repository uses [pre-commit](https://pre-commit.com/) to run some checks on the code before committing. You can install it with `pip install pre-commit` and then run `pre-commit install` to install the git hooks.

Thank you!
