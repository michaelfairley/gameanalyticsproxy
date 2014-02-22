package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		hasher := md5.New()
		hasher.Write(body)
		io.WriteString(hasher, os.Getenv("GAME_ANALYTICS_SECRET_KEY"))

		hexdigest := hex.EncodeToString(hasher.Sum(nil))

		request, err := http.NewRequest("POST", "http://api.gameanalytics.com"+r.URL.Path, bytes.NewReader(body))
		if err != nil {
			panic(err)
		}
		request.Header.Set("Authorization", hexdigest)

		client := &http.Client{}
		resp, err := client.Do(request)
		fmt.Println(resp.Status)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(204)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
