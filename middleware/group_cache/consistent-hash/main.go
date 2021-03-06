package main

import (
	"fmt"
	"log"
	"net/http"

	single_node "gowhole/middleware/group_cache/http-server/single-node"
)

/*
   $ curl http://localhost:9999/_geecache/scores/Tom
   630
   $ curl http://localhost:9999/_geecache/scores/kkk
   kkk not exist
*/

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	single_node.NewGroup("scores", 2<<10, single_node.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := single_node.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
