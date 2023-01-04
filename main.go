package main

import (
	"log"
	"net"
	"net/http"
)

// //go:embed www
// var fs embed.FS

func main() {
	authSpotify()

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

	// Handle requests on /ai/answer with Send function
	http.HandleFunc("/ai/answer", Send)

	//serv on port 5019
	http.Serve(ln, nil)

	log.Println("exiting...")
}
