package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nheingit/learnGo/cli"
)

//--
func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// Just put some headers to allow CORS...
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}

//--

func logic(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	bodyData := buf.String()

	log.Print("bodyData : ")
	log.Println(bodyData)

	list := strings.Split(bodyData, ",")
	fmt.Println(list)

	inputls := []string{}

	for i := len(list) - 1; i >= 0; i-- {
		data := strings.Replace(list[i], "[", "", -1)
		data = strings.Replace(data, "]", "", -1)
		data = strings.Replace(data, "\"", "", -1)
		inputls = append(inputls, data)
	}

	log.Println("inputls")
	log.Println(inputls)

	gob.Register(elliptic.P256())
	//defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Populate()

	outputls := cmd.CheckPopulation(inputls)
	fmt.Print(outputls)

	jsonBytes, _ := json.Marshal(outputls)
	io.WriteString(response, string(jsonBytes))

}
func manejadorSolicitudes() {
	//enrutador
	r := mux.NewRouter()
	enableCORS(r)
	//endpoints
	r.HandleFunc("/logic", logic)
	log.Fatal(http.ListenAndServe(":9001", r))

}

func main() {

	manejadorSolicitudes()

}
