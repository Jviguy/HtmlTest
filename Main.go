package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := os.Mkdir("./public", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Println(err)
		_, _ = fmt.Scan()
		return
	}
	fs := interceptHandler(http.FileServer(http.Dir("./public")),Error)
	http.Handle("/",fs)
	fmt.Println("Starting webserver on port 3000")
	err = http.ListenAndServe(":3000",nil)
	if err != nil{
		fmt.Println("Couldn't start webserver on port 3000 counting up")
		for i := 3000; i > 4000;i++ {
			err = http.ListenAndServe(fmt.Sprint(":",i),nil)
			if err != nil{
				fmt.Println("Starting webserver on port " + strconv.Itoa(i))
				break
			}
		}
	}
}
type interceptResponseWriter struct {
	http.ResponseWriter
	*http.Request
	header http.Header
	errH func(http.ResponseWriter,*http.Request, int)
}

func (w *interceptResponseWriter) Header() http.Header {
	return w.header
}

func (w *interceptResponseWriter) WriteHeader(status int) {
	if status >= http.StatusBadRequest {
		w.errH(w.ResponseWriter,w.Request, status)
		w.errH = nil
	} else {
		w.ResponseWriter.WriteHeader(status)
	}
}

type ErrorHandler func(http.ResponseWriter, *http.Request,int)

func (w *interceptResponseWriter) Write(p []byte) (n int, err error) {
	if w.errH == nil {
		return len(p), nil
	}
	return w.ResponseWriter.Write(p)
}

func defaultErrorHandler(w http.ResponseWriter,r *http.Request,status int) {
	log.Print("error handler called")
	http.Error(w, "foo", status)
}

func interceptHandler(next http.Handler, errH ErrorHandler) http.Handler {
	if errH == nil {
		errH = defaultErrorHandler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&interceptResponseWriter{w, r, w.Header() ,errH}, r)
	})
}

func Error(w http.ResponseWriter,r *http.Request, status int)  {
	if status != 200 {
		http.ServeFile(w, r, fmt.Sprint("./public/", status, ".html"))
	}
}