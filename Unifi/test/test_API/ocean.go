package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func mainNOT() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")

	/*
		second := r.URL.Query().Get("second")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
		}

		fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s, body:\n%s\n",
			ctx.Value(keyServerAddr),
			hasFirst, first,
			hasSecond, second,
			body)
		io.WriteString(w, "This is my website!\n")
	*/

}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
