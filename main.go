package main

import (
	"flag"
	"io/ioutil"
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

func ReportIP() {
	resp, err := http.Get(`http://ip.clearcode.cn/ip`)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("ip", string(data))
}

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

var (
	isReporter bool
)

func init() {
	flag.BoolVar(&isReporter, "r", false, "is reporter")
}

func main() {
	flag.Parse()

	if isReporter {
		tk := time.NewTicker(30 * time.Second)
		for range tk.C {
			ReportIP()
		}
	} else {
		StartHttpServer()
	}
}
