package balancing

import (
	"net/http"
	"net/rpc"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
)

func login(body []byte, client *rpc.Client) ([]byte, error) {
	var resp interface{}
	var balancing Balancing

	err := client.Call("Remote.Login", body, &resp)
	if err != nil {
		log.Println(err.Error(), "xxxxxxxxxxxxxx")
		balancing.Code = "500"
		balancing.Msg = "服务器内部错误. " + err.Error()
		return []byte{}, &balancing
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error(), "yyyyyyyyyyyyyy")
		balancing.Code = "500"
		balancing.Msg = "服务器内部错误. " + err.Error()
		return []byte{}, &balancing
	}

	return b, nil
}

func user(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var index int64
	index = 0
	h := (map[string][]string(r.Header))
	org, ok := h["org"]
	if ok {
		var err1 error
		index , err1 = strconv.ParseInt(org[0], 10, 64)
		if err1 != nil {
			log.Println(err1.Error())
			w.Write([]byte(`{"code":400,"msg":"` + err1.Error() + `"}`))
			return
		}
	}

	switch r.Method {
	case "POST":
		client := RPCConn(index)
		defer client.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(`{"code":400,"msg":"` + err.Error() + `"}`))
			return
		}

		log.Println(r.URL.Path)
		switch r.URL.Path {
		case "/user/login":
			b, err := login(body, client)
			if err != nil {
				log.Println(err.Error())
				w.Write([]byte(`{"code":"` + err.(*Balancing).Code + `","msg":"` + err.Error() + `"}`))
				return
			}
			w.Write(b)
		}
	default:
		w.Write([]byte(`{"code":404,"msg":"Page not found"}`))
	}
}