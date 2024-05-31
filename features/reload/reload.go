package reload

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/fsnotify/fsnotify"
)

func Handler(a *appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create new watcher.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		// Start listening for events.
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					fmt.Fprintf(w, fmt.Sprintf("data: %s\n\n", event.Name))
					w.(http.Flusher).Flush()
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					a.Log.Error("reload watcher", "error", err)
				}
			}
		}()

		// Add a path.
		err = watcher.Add("content/pages")
		if err != nil {
			log.Fatal(err)
		}

		<-r.Context().Done()
	}
}
