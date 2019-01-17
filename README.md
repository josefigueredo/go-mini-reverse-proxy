# go-mini-reverse-proxy
Situation: I only have one port open in a server and I need to consume the web ui that was open in another. I didn't have sudo so I build this little tool.

Compilation: go build -o mini-proxy main.go

Execution: ./mini-proxy -listeningPort=":8032" -forwardingPort=":8088"
