# telegram_bot

## Problem

Running multiple scripts with telegram reporting using the same bot could result in conflicts. This tool provides HTTP API to send messages.

## Usage

1. Specify telegram token in ```config.yaml```:

    ```yml
    telegram_token: "token goes here"
    ```

2. Run ```telegram_bot```. See ```telegram_bot --help``` for command line options
3. Get chat ID with one of two ways
    1. Start conversation, send message to bot mentioning it
    2. Add your bot to a group. It should report group id now. To get ID of a group if bot is already a member [send a message that starts with `/`](https://core.telegram.org/bots#privacy-mode)

### Sending messages

```sh
curl -d 'test' "127.0.0.1:9031/$chat_id?mode=Markdown"
curl -d @message.txt "127.0.0.1:9031/$chat_id"
```

Set ```chat_id``` to the number you got from your bot, with ```-``` if it was reported so (true for groups).

## Test your instance

1. Build
2. Export TELEGRAM_CHATID environment variable
3. Run `prove`

```bash
go build
export TELEGRAM_CHATID="-YOUR TELEGRAM CHAT ID"
prove -v
```
