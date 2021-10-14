// create server and return a random 5 char long text as response
package utils
import (
	"fmt"
	"strings"
	"net/http"
)

func response(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, strings.Repeat("hello", 100))
}

func StartServer() () {
	http.HandleFunc("/", response)
	http.ListenAndServe(":10427", nil)
}