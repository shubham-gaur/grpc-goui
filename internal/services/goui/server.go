package goui

import (
	"bytes"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var SupportedPorts = flag.String("ports", os.Getenv("EXPOSED_PORTS"), "List of ports to be exposed")
var SupportedIntf = flag.String("intf", os.Getenv("INTERFACE"), "Server external interface name")

type UIServer struct {
	Cmds  Cmds
	Ports map[int]string
}

type Service interface {
	PlaintextServer(addrs string, port string) (string, Cmds, map[int]string)
	KillServer(uri string, port int) map[int]string
}

type uiserver struct {
	server   string
	uiserver *UIServer
}
type Cmds map[string]Process

type Process struct {
	Cmd *exec.Cmd
	Pid int
}

func NewUIserver(s UIServer) Service {
	return &uiserver{uiserver: &s}
}

func (g *uiserver) PlaintextServer(addrs string, port string) (string, Cmds, map[int]string) {
	flag.Parse()
	log.Println("Spawning Server for " + addrs + ":" + port)
	log.Println("Executing gRPC UI ...")
	supportedPorts := strings.Split(*SupportedPorts, "...")
	portMin, _ := strconv.Atoi(supportedPorts[0])
	portMax, _ := strconv.Atoi(supportedPorts[1])
	var portAssigned int
	var cmd *exec.Cmd
	buf := make([]byte, 1024)
	cmp := make([]byte, 1024)
	for portAssigned = portMin; portAssigned < portMax; portAssigned++ {
		if g.uiserver.Ports == nil {
			g.uiserver.Ports = make(map[int]string)
			g.uiserver.Ports[portAssigned] = addrs + ":" + port
			cmd, buf = startGrpcUI(portAssigned, addrs, port)
			break
		} else if _, ok := g.uiserver.Ports[portAssigned]; !ok {
			g.uiserver.Ports[portAssigned] = addrs + ":" + port
			cmd, buf = startGrpcUI(portAssigned, addrs, port)
			break
		} else {
			log.Println("No free ports available!")
		}
	}

	if !bytes.Equal(buf, cmp) {
		out := string(buf[0:100])
		log.Println("Command o/p generated " + out)
		outArr := strings.Split(out, " ")
		uri := outArr[len(outArr)-1]
		uriFmt := strings.Split(uri, "/")
		g.server = uriFmt[0] + "//" + uriFmt[2]
		log.Println("Server started on " + uri)
		if g.uiserver.Cmds == nil {
			g.uiserver.Cmds = make(Cmds)
		}
		g.uiserver.Cmds[g.server] = Process{
			Cmd: cmd,
			Pid: cmd.Process.Pid,
		}
		log.Println("Process started " + strconv.Itoa(g.uiserver.Cmds[g.server].Pid))
	} else {
		pid := cmd.Process.Pid
		cmd.Process.Wait()
		log.Println("Process stopped " + strconv.Itoa(pid))
	}
	return g.server, g.uiserver.Cmds, g.uiserver.Ports
}

func (g *uiserver) KillServer(uri string, port int) map[int]string {
	log.Println("Killing Pid..." + strconv.Itoa(g.uiserver.Cmds[uri].Pid))
	delete(g.uiserver.Ports, port)
	if err := g.uiserver.Cmds[uri].Cmd.Process.Kill(); err != nil {
		log.Fatal("Failed to kill process: ", err)
	}
	return g.uiserver.Ports
}

func startGrpcUI(portAssigned int, addrs string, port string) (*exec.Cmd, []byte) {
	flag.Parse()
	buf := make([]byte, 1024)
	var cmd *exec.Cmd
	extInterface := *SupportedIntf
	log.Println("Interface chosen " + extInterface)
	log.Println("Fetching IP address...")
	localAddrs, _ := getIPFromINTF(extInterface)
	log.Println("Fetched local address " + localAddrs)
	log.Println("Port assigned " + strconv.Itoa(portAssigned))

	log.Println("Command to be executed " + "grpcui" + " " + "-open-browser=false" + " " + "-plaintext" + " " + "-port=" + strconv.Itoa(portAssigned) + "-bind=" + localAddrs + " " + addrs + ":" + port)
	cmd = exec.Command("grpcui", "-open-browser=false", "-plaintext", "-port="+strconv.Itoa(portAssigned), "-bind="+localAddrs, addrs+":"+port)
	log.Println("Command executed ...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicln(err)
	}
	go func() {
		if err := cmd.Start(); err != nil {
			log.Panicln(err)
		}
	}()
	stdout.Read(buf)
	return cmd, buf
}
func getIPFromINTF(intf string) (string, error) {
	intfName, err := net.InterfaceByName(intf)
	if err != nil {
		log.Println(err)
	}
	log.Println(intf)
	log.Println(intfName)
	log.Println("Interface name " + intfName.Name)
	addrs, err := intfName.Addrs()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Addresses on interface %v are %v", intfName.Name, addrs)
	ipAddrs := addrs[0].String()
	ipAddrsWithoutSubnet := strings.Split(ipAddrs, "/")
	return ipAddrsWithoutSubnet[0], err
}
