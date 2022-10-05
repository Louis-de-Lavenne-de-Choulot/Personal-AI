package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/fstanis/screenresolution"
	"github.com/zserge/lorca"
)

// //go:embed www
// var fs embed.FS

func main() {
	authSpotify()
	args := []string{}
	resolution := screenresolution.GetPrimary()
	if resolution == nil {
		fmt.Println("failed to get screen resolution")
		os.Exit(1)
	}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", resolution.Width, resolution.Height, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	mod := &Mod{}
	ui.Bind("Ans", mod.Send)
	ui.Bind("Send", mod.Send)
	ui.Bind("GetInput", mod.GetInput)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:5019")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	// serv /www/ as root page
	http.Handle("/", http.FileServer(http.Dir("./www")))
	// //redirect /callback to root page
	// http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// })
	//serv on port 5019
	go http.Serve(ln, nil)

	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
