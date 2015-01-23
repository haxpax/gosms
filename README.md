gosms
-----

Your own local SMS gateway

- deploy in less than 1 minute
- supports Windows, GNU\Linux, Mac OS
- works with GSM modems
- provides API over HTTP to push messages to gatewy, just like the internet based gateways do
- takes care of queing, throttling, retrying etc

![gosms dashboard](https://github.com/haxpax/gosms/blob/screenshot/screenshots/gosms.png)

deployment
----------
- download suitable binary from downloads section
- create sqlite database using misc/queries.sql
- edit conf.ini to match your configuration
- execute binary

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
      - 
        ```json
        {
          "status": 200,
          "message": "ok"
        }
        ```
- /api/logs/ [*GET*]
    - response
        - 
          ```json
          {
            "status": 200,
            "message": "ok",
            "summary": [ <total_pending>int, <total_processed>int, <total_error>int ],
            "daycount": { <date>string: <count>int, },
            "messages": [
              {
                "uuid": string,
                "mobile": string,
                "body": string,
                "status": int
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