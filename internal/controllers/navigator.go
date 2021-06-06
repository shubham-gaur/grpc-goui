package controllers

import (
	"log"
	"strconv"
	"strings"

	"github.com/shubham-gaur/goui/internal/services/goui"

	"github.com/gin-gonic/gin"
)

var (
	Navigator      = navigator{}
	cmds           goui.Cmds
	servers        map[string]string
	gRPC           string
	SupportedPorts map[int]string
)

type navigator struct{}

func (controllers navigator) Index(ctx *gin.Context) {
	ctx.HTML(200, "index", gin.H{
		"Title": "Test",
	})
}

func (controllers navigator) Server(ctx *gin.Context) {
	url, _ := ctx.GetPostForm("server_id")
	port, _ := ctx.GetPostForm("server_port")
	addrs := url + ":" + port
	if servers == nil {
		servers = make(map[string]string)
	}
	log.Println("gRPC Sever Info " + url + ":" + port)
	_, ok := servers[addrs]
	if !ok {
		uiserver := goui.NewUIserver(goui.UIServer{Cmds: cmds, Ports: SupportedPorts})
		gRPC, cmds, SupportedPorts = uiserver.PlaintextServer(url, port)
		pid := strconv.Itoa(cmds[gRPC].Pid)
		log.Println("gRPC UI " + gRPC)
		log.Println("gRPC Pid " + pid)
		if gRPC != " " && addrs != ":" && cmds[gRPC].Pid != 0 {
			servers[addrs] = gRPC
		}
		ctx.HTML(200, "index", gin.H{
			"Title":   "Test Success",
			"Servers": servers,
		})
	} else {
		ctx.HTML(200, "index", gin.H{
			"Title":   "Server already running...",
			"Servers": servers,
		})
	}
}

func (controllers navigator) KillServer(ctx *gin.Context) {
	server_id := ctx.Param("serverid")
	url, ok := servers[server_id]
	urlFmt := strings.Split(url, ":")
	port, _ := strconv.Atoi(urlFmt[2])
	if ok {
		log.Println("server_id received ", server_id)
		log.Println("Killing gRPC UI " + url)
		uiserver := goui.NewUIserver(goui.UIServer{Cmds: cmds, Ports: SupportedPorts})
		SupportedPorts = uiserver.KillServer(url, port)
		log.Println("Supported ports updated... ")
		log.Println(SupportedPorts)
		delete(servers, server_id)
	}
	ctx.HTML(200, "index", gin.H{
		"Title":   "Test Success",
		"Servers": servers,
	})
}
