package main

import (
	"fmt"
	"net/http"
	"os"
)

const PORT = "8081"
func main(){
	mux := new(http.ServeMux)
	server := http.Server{
		Addr: ":" + PORT,	
		Handler: mux,
	}

	fmt.Printf("Server running on port: %s\n", PORT)

	err := server.ListenAndServe()
	if err != nil{
		fmt.Print(err)
		os.Exit(1)
	}

}
