// superproxy project superproxy.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/CreditTone/aslist"
	colorfulog "github.com/CreditTone/colorfulog"
	"github.com/lxzan/hasaki"
)

var (
	//你的代理提取api
	XUN_PROXY_API_URL = "http://api.xdaili.cn/xdaili-api//privateProxy/applyStaticProxy?spiderId=xxx&returnType=2&count=1"
	proxyList         = aslist.NewAsList()
)

func setGlobalProxyList(pl []interface{}) {
	if len(pl) == 0 {
		return
	}
	for _, p := range pl {
		m := p.(map[string]interface{})
		proxy := fmt.Sprintf("%s:%v", m["ip"], m["port"])
		colorfulog.Infof("add proxy %v", proxy)
		proxyList.LeftPush(proxy)
	}
	for proxyList.Length() > 75 {
		colorfulog.Infof("remove proxy %v", proxyList.RightPop())
	}
}

func updateUpstreamProxy() {
	defer func() {
		if err := recover(); err != nil {
			colorfulog.Warnf("fetchUpstreamProxy %v", err)
		}
	}()
	body, err := hasaki.Get(XUN_PROXY_API_URL).Send(nil).GetBody()
	if err != nil {
		return
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal(body, &mapResult)
	if err != nil {
		colorfulog.Warnf("JsonToMapDemo err: %v ", err)
		return
	}
	ERRORCODE := mapResult["ERRORCODE"]
	colorfulog.Infof("fetchUpstream errorcode %s", ERRORCODE)
	if ERRORCODE == "0" {
		result := mapResult["RESULT"]
		if proxyList, ok := result.([]interface{}); ok {
			colorfulog.Infof("add proxyList : %v", proxyList)
			setGlobalProxyList(proxyList)
		} else {
			colorfulog.Warnf("convert field %v", result)
		}
	}
}

func forward(conn net.Conn, remoteAddr string) {
	client, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		defer conn.Close()
		return
	}
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), client.RemoteAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func updateProxyList() {
	colorfulog.Info("启动代理更新协程")
	for {
		updateUpstreamProxy()
		time.Sleep(time.Second * 20)
	}
}

func main() {

	//启动代理更新线程
	go updateProxyList()

	listener, err := net.Listen("tcp", ":12001")
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		var px interface{}
		px = proxyList.RandomGet()
		if px == nil {
			//随便设置个无效的代理
			px = "127.0.0.1:60001"
		}
		go forward(conn, px.(string))
	}
}
