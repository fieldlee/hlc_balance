// balancing

package balancing

import (
	"io/ioutil"
	"encoding/json"
	"net/http"

	"log"
)

type Server struct {
	Name	string	`json:"name"`
	Domain_port	string	`json:"domain_port"`
}

type Config struct {
	Listen	string	`json:"listen"`
	Servers	[]Server	`json:"servers"`
	Global	[]Server	`json:"global"`
}

var config Config

type Balancing struct {
	Code	string	`json:"code"`
	Msg	string	`json:"msg"`
	Data	interface{}	`json:"data"`
}

func (this Balancing)Error() string {
	return this.Msg
}

func Run() {
	data, err := ioutil.ReadFile("/etc/hlc/hlc-blc.conf.json")
	if err != nil {
		log.Fatalln(err.Error(), "cannot find the file: balancing.conf.json")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(err.Error(), "cannot parse the file: config.json")
	}

	http.HandleFunc("/user/login", user)
	http.HandleFunc("/asset/register", asset)
	http.HandleFunc("/asset/querydetail", asset)
	http.HandleFunc("/asset/feed", asset)
	http.HandleFunc("/asset/medication", asset)	// 检疫
	http.HandleFunc("/asset/prevention", asset)	// 防疫
	http.HandleFunc("/asset/save", asset)
	http.HandleFunc("/asset/lost", asset)
	http.HandleFunc("/asset/fattened", asset)

	log.Println("Balancing...", config.Listen)

	log.Fatalln(http.ListenAndServe(config.Listen, nil))
}
