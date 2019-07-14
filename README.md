
### :floppy_disk: Bank

###### install
`docker-compose up`

###### usage
```
curl -X POST '0.0.0.0:8080' \
-d '{"key":"key", "value":"value", "avail":[{"start":"3:59am","duration":30}]}' \
-H "Content-Type: application/json"
```
