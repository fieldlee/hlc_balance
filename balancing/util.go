// balancing

package balancing

import (
	"time"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
	"os"
	"strings"
)

func RPCConn(index int64) *rpc.Client {
	//index := time.Now().Unix() % int64(len(config.Servers))

	client, err := jsonrpc.Dial("tcp", config.Servers[index].Domain_port)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(config.Servers[index].Domain_port)

	return client
}

func ConnectGlobalServer(index int64) *rpc.Client {
	log.Println(config.Servers[index].Domain_port)
	global, err := jsonrpc.Dial("tcp", config.Global[index].Domain_port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return global
}

//打印内容到文件中
//tracefile(fmt.Sprintf("receive:%s",v))
func Tracefile(filename string,str_content string)  {
	fd,_:=os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	fd_time:=time.Now().Format("2006-01-02 15:04:05");
	fd_content:=strings.Join([]string{"======",fd_time,"=====",str_content,"\n"},"")
	buf:=[]byte(fd_content)
	fd.Write(buf)
	fd.Close()
}
