package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	kv    = make(map[string]string)
)

func StartHttpServer() {
	http.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		v := r.RemoteAddr
		host, _, _ := net.SplitHostPort(v)
		v = host

		ip := r.Header.Get("X-Real-IP")
		if ip != "" {
			log.Println(ip)
			v = ip
		}

		mutex.Lock()
		kv[key] = v
		mutex.Unlock()
		w.Write([]byte(v))
	})

	http.HandleFunc("/kv", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		mutex.Lock()
		defer mutex.Unlock()

		v := kv[key]
		w.Write([]byte(v))
	})

	http.ListenAndServe(":9999", nil)
}

func startTicker() {
	fn := func() {
		addrs, err := net.LookupHost("mrjnamei.tpddns.cn")
		if err != nil {
			log.Println("[ERRO] ", err)
			return
		}
		if len(addrs) == 0 {
			return
		}
		mutex.Lock()
		kv["ip"] = addrs[0]
		mutex.Unlock()
		log.Println("[INFO] ip is", addrs[0])
	}

	fn()
	tk := time.NewTicker(30 * time.Second)
	for range tk.C {
		fn()
	}
}

func main() {
	go startTicker()
	StartHttpServer()
}
