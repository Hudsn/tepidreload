package tepidreload

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed iframe.tmpl
var iframeFS embed.FS

//go:embed script.tmpl
var scriptFS embed.FS

// Returns two handlers: One to serve a generated script file to poll for reloads, and the second is the endpoint that the polling script checks against.
func MakeHandlers(checkPath string, listenPort string, config Config) (http.HandlerFunc, http.HandlerFunc) {
	t, err := template.ParseFS(scriptFS, "*.tmpl")
	if err != nil {
		panic(err)
	}

	type Data struct {
		Path string
		Port string
	}

	data := Data{
		Path: checkPath,
		Port: listenPort,
	}

	socketHandler := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
			return
		}

		writer := makeWritefunc(config)

		go writer(ws)
	}

	scriptHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		t.ExecuteTemplate(w, "script.tmpl", data)
	}

	return scriptHandler, socketHandler
}

func IframeHandler() http.HandlerFunc {
	type TemplateData struct {
		StaticPath string
	}

	t, err := template.ParseFiles("iframe.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		data := TemplateData{
			StaticPath: r.URL.Path,
		}

		t.ExecuteTemplate(w, t.Name(), data)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func makeWritefunc(config Config) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		fileTicker := time.NewTicker(time.Millisecond * time.Duration(config.TickIntervalMS))
		defer func() {
			fileTicker.Stop()
			ws.Close()
		}()

		for {
			select {
			case <-fileTicker.C:
				isChanged, _ := checkFileMods(config)
				if isChanged {
					if err := ws.WriteMessage(websocket.TextMessage, []byte("true")); err != nil {
						return
					}
				}
			}
		}

		//
		//
	}
}
