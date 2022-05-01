package main

import (
	"fmt"
	"net/http"
	"time"
)

func foo(w http.ResponseWriter, req *http.Request) {
	fmt.Println(time.Now(), "foo called")
	//	time.Sleep(4 * time.Second)
	fmt.Fprintf(w, "foo\n")
	fmt.Fprintf(w, "bar\n")
	fmt.Fprintf(w, "baz\n")
	fmt.Fprintf(w, "!\n")
}

func bar(w http.ResponseWriter, req *http.Request) {
	fmt.Println(time.Now(), "bar called")
	time.Sleep(4 * time.Second)
	for name, vals := range req.Header {
		for _, val := range vals {
			fmt.Fprintf(w, "%v: %v\n", name, val)
		}
	}
}

func main() {
	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
