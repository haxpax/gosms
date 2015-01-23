gosms
-----

Your own local SMS gateway

- deploy in less than 1 minute
- supports Windows, GNU\Linux, Mac OS
- works with GSM modems
- provides API over HTTP to push messages to gateway, just like the internet based gateways do
- takes care of queuing, throttling and retrying

![gosms dashboard](https://github.com/haxpax/gosms/blob/screenshot/screenshots/gosms.png)

deployment
----------
- Update `conf.ini` [DEVICES] section with your modem's COM port (for ex. COM10 or /dev/USBtty2)
- Run

API specification
------------------
- /api/sms/ [*POST*]
    - param **mobile**
        - mobile number to send message to
        - number should have contry code prefix
        - for ex. +919890098900
    - param **message**
        - message text
        - max length is limited to 160 characters
    - response
```json
{
  "status": 200,
  "message": "ok"
}
```
- /api/logs/ [*GET*]
    - response
```json
{
  "status": 200,
  "message": "ok",
  "summary": [ 10, 50, 2 ],
  "daycount": { "2015-01-22": 10, "2015-01-23": 25 },
  "messages": [
    {
      "uuid": "d04f17c4-a32c-11e4-827f-00ffcf62442b",
      "mobile": "+1858111222",
      "body": "Hey! Just playing around with gosms.",
      "status": 1
    },
  ]
}
```
    - message status codes
      - 0 : Pending
      - 1 : Processed
      - 2 : Error

planned features
-------
- CRUD support for messages, possibly support cancellation of message
- authentication support for API
- authentication support for WebUI
- support multiple devices with load balancing


building on windows
-------------------
- go get `github.com/haxpax/gosms`
- Setup GCC for go-sqlite3 package
	- Download MinGW from http://sourceforge.net/projects/mingw/
	- Add `C:\MinGW\bin` to PATH
	- run `mingw-get install gcc` from command line
- go build