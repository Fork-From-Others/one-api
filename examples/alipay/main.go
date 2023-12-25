package main

import "net/http"

func main() {
	http.HandleFunc("/alipay", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello alipay"))
		if err != nil {
			return
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
