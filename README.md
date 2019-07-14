
# Bank

###### install
`docker-compose up`

###### usage
```
POST http://0.0.0.0:8080
Content-Type: application/json

{
    "key":"key",
    "value":"value",
    "avail":[
        { "start":"4:00pm","duration":30 }  // server time, duration in seconds
    ]
}
```

```
GET http://0.0.0.0:8080/key
```

###### environment variables
```
OUTPUT=output/file            // set output file.
DEBUG=false                   // set output to visible.

ALLOWED=*                     // all.
ALLOWED=127.0.0.1             // single ip.
ALLOWED=127.0.0.1,127.0.0.1   // multiple ip.
```
