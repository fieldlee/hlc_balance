// balancing
package balancing

import (
	"net/http"
	"net/rpc"
	"io/ioutil"
	"encoding/json"
	"log"
)

func sessionValidation(header map[string][]string, global *rpc.Client) error {
	var resp bool
	var balancing Balancing

	err := global.Call("Remote.SessionValidation", header, &resp)
	if err != nil {
		balancing.Code = "500"
		balancing.Msg = "服务器内部错误. " + err.Error()
		return &balancing
	}

	if !resp {
		balancing.Msg = "403"
		balancing.Msg = "wrong token"
		return &balancing
	}

	return nil
}

func assetFunc(header map[string][]string, body []byte, remoteFunc string, client *rpc.Client) ([]byte, error) {
	var resp map[string]interface{}
	var balancing Balancing

	args := make(map[string]map[string][]string)
	args["header"] = header
	args["body"] = make(map[string][]string)
	args["body"]["b"] = []string{string(body)}

	err := client.Call(remoteFunc, args, &resp)
	if err != nil {
		balancing.Code = "500"
		balancing.Msg = "服务器内部错误. " + err.Error()
		return []byte{}, &balancing
	}

	balancing.Code = "0"
	balancing.Msg = "success"
	balancing.Data = resp
	respJSON, err := json.Marshal(balancing)
	if err != nil {
		balancing.Code = "500"
		balancing.Msg = "服务器内部错误. " + err.Error()
		return []byte{}, &balancing
	}

	return respJSON, nil
}

func asset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	global := ConnectGlobalServer(config.Global.Domain_port)
	defer global.Close()

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(`{"code":400,"msg":"` + err.Error() + `"}`))
			return
		}

		err = sessionValidation(map[string][]string(r.Header), global)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(`{"code":"` + err.(*Balancing).Code + `","msg":"` + err.Error() + `"}`))
			return
		}
		client := RPCConn()
		defer client.Close()

		var respJSON []byte
		switch r.URL.Path {
		case "/asset/register":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetRegister", client)
		case "/asset/querydetail":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetQueryDetail", client)
		case "/asset/feed":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetFeed", client)
		case "/asset/medication":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetMedication", client)
		case "/asset/prevention":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetPrevention", client)
		case "/asset/save":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetSave", client)
		case "/asset/lost":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetLost", client)
		case "/asset/fattened":
			respJSON, err = assetFunc(map[string][]string(r.Header), body, "Remote.AssetFattened", client)
		default:
			w.Write([]byte(`{"code":404,"msg":"Page not found"}`))
			return
		}

		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(`{"code":"` + err.(*Balancing).Code + `","msg":"` + err.Error() + `"}`))
			return
		}
		w.Write(respJSON)

		return
	default:
		w.Write([]byte(`{"code":404,"msg":"Page not found"}`))
	}
}
