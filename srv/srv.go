package main

import (
	"demo/pb_ws_test/srv/control"
	"flag"
	log "github.com/alecthomas/log4go"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }} // use default options

func wsEcho(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}

	go func() {
		defer func() {
			log.Error("exit")
			err := c.Close()
			if err != nil {

			}
		}()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Error(err)
				break
			}

			if mt == websocket.TextMessage {
				err := c.WriteMessage(mt, message)
				if err != nil {
				}

			} else {
				err := control.GetHub().DispatchMessage(c, message)
				if err != nil {
				}
			}
		}
	}()
}

func home(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
	if err != nil {

	}
}

func main() {
	log.LoadConfiguration("../configs/log4go.xml")

	flag.Parse()

	http.HandleFunc("/ws", wsEcho)

	http.HandleFunc("/", home)

	log.Error(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8" />
	<script>  
		window.addEventListener("load", function(evt) {
	
	    var output = document.getElementById("output");
	    var input = document.getElementById("input");
	    var ws;
	
	    var print = function(message) {
	        var d = document.createElement("div");
	        d.innerHTML = message;
	        output.appendChild(d);
	    };
	
	    document.getElementById("open").onclick = function(evt) {
	        if (ws) {
	            return false;
	        }
	        ws = new WebSocket("{{.}}");
	        ws.onopen = function(evt) {
	            print("OPEN");
	        }
	        ws.onclose = function(evt) {
	            print("CLOSE");
	            ws = null;
	        }
	        ws.onmessage = function(evt) {
	            print("RESPONSE: " + evt.data);
	        }
	        ws.onerror = function(evt) {
	            print("ERROR: " + evt.data);
	        }
	        return false;
	    };
	
	    document.getElementById("send").onclick = function(evt) {
	        if (!ws) {
	            return false;
	        }
	        print("SEND: " + input.value);
	        ws.send(input.value);
	        return false;
	    };
	
	    document.getElementById("close").onclick = function(evt) {
	        if (!ws) {
	            return false;
	        }
	        ws.close();
	        return false;
	    };
	
	});
	</script>
	</head>
	<body>
	<table>
	<tr><td valign="top" width="50%">
	<p>Click "Open" to create a connection to the server, 
	"Send" to send a message to the server and "Close" to close the connection. 
	You can change the message and send multiple times.
	<p>
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="input" type="text" value="Hello world!">
	<button id="send">Send</button>
	</form>
	</td><td valign="top" width="50%">
	<div id="output"></div>
	</td></tr></table>
	</body>
	</html>
`))
