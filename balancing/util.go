// balancing

package balancing

import (
	"time"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
)

func RPCConn() *rpc.Client {
	index := time.Now().Unix() % int64(len(config.Servers))

	for {
		client, err := jsonrpc.Dial("tcp", config.Servers[index].Domain_port)
		if err != nil {
			index++
			if int(index) == len(config.Servers) {
				index ^= index
			}
			continue

			log.Println(err.Error())
		}
		log.Println(config.Servers[index].Domain_port)

		return client
	}
}

func ConnectGlobalServer(domain_port string) *rpc.Client {
	global, err := jsonrpc.Dial("tcp", domain_port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return global
}
