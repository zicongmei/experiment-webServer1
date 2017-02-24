package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "Hello web!\n")
	fmt.Fprintf(w, "Host = %v\n", r.URL.Host)
	fmt.Fprintf(w, "RawPath = %v\n", r.URL.RawPath)
	fmt.Fprintf(w, "Path = %v\n", r.URL.Path)
	fmt.Fprintf(w, "key = %v\n", r.Form.Get("key"))
}

func tokenSigninHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	token := r.Form.Get("token")
	log.Print("token=" + token)

	endPoint := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
	if response, err := http.Get(endPoint); err != nil{
		fmt.Fprintln(w, "Error in checking endpoint " + err.Error())
	} else {
		fmt.Fprintln(w, "Endpoint return code: " + response.Status)
		if response.StatusCode == 200 {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			contentMap := make(map[string]string)
			if err := json.Unmarshal(body, &contentMap); err != nil {
				fmt.Fprintln(w, "can't unmarshall body")
			} else {
				keys := []string{"name", "email"}
				for _, k := range keys {
					if val, exist := contentMap[k]; exist {
						fmt.Fprintln(w, k + " = " + val)
					}
				}
			}
			fmt.Fprintf(w, "===========\nBody:\n%v\n", string(body))
		} else {
			fmt.Fprintln(w, "token invalid")
		}
	}
}

func main() {
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("login"))))
	http.HandleFunc("/token-signin/", tokenSigninHandler)
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
