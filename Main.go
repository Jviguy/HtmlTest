package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := os.Mkdir("./public", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Println(err)
		_, _ = fmt.Scan()
		return
	}
	http.FileServer(http.Dir("/"))
}
