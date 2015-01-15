gosms
-----

Your own local SMS gateway

- deploy in less than 1 minute
- supports Windows, GNU\Linux, Mac OS
- works with GSM modems
- provides API over HTTP to push messages to gatewy, just like the internet based gateways do
- takes care of queing, throttling, retrying etc

deployment
----------
- download suitable binary from downloads section
- create sqlite database using misc/queries.sql
- edit conf.ini to match your configuration
- execute binary

API specification
------------------
- /api/sms
    - POST
    - param **mobile**
        - mobile number to send message to
        - Number should have contrycode prefix
        - for example: +919890098900
    - param **message**
        - message text
        - max length is limited to 160 characters
    - response
        - status 200, "ok" on success

- /api/smsdata/<start>
    - GET
    - **start** should be an integer specifying starting offset
    - response
        - ```json
            {   "messages": [   
                        {   "uuid": string,
                            "mobile": string,
                            "body": string,
                            "status": int
                        },
                    ]
                "iDisplayStart": int,
                "iDisplayLength": int,
                "iTotalRecords": int,
                "iTotalDisplayRecords": int
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