# screamery

### a chat app

go / htmx / bbolt / bulma.

## features
- local authentication
- create posts for your profile
- comment on other users posts
- message sanitization
- local db (bbolt)
- messages are only stored in memory (never written to disk)
- send "commands" to the server (see help page)
- a help page
- stats and graphs, if you're into that kind of thing


## install and run
`go build . # in root directory`

```
./scream -h
Usage of ./scream:
  -cert-file string
    	cert file (default "server-cert.pem")
  -chart-service-log string
    	chart service log (default "charting_service.log")
  -chart-service-port string
    	chart service port (default ":10440")
  -chart-service-url string
    	chart service url (default "http://localhost:10440/graph")
  -db-name string
    	database name (default "chat.db")
  -first-user-mode
    	first user mode
  -key-file string
    	key file (default "server-key.pem")
  -log-file string
    	log file (default "chat.log")
  -message-limit int
    	message limit (default 100)
  -self-host
    	self host microservice (default true)
  -token-bucket string
    	token bucket (default "tokens")
  -update-freq duration
    	update frequency (default 2m0s)
  -url string
    	url (default ":8081")
  -user-bucket string
    	user bucket (default "users")
```

### endpoints
- localhost:8080/add-user to add first user.
- /help for help
- / to hit login
- /room/whatever to create a new room
- /stats to see the charts

## example

![room_example](docs/example.png)
