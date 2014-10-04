package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"

	"protocol"
	"tunnel"
	"userdb"

	"github.com/hamo/golog"
)

var (
	sc  *serverConfig
	udb *userdb.UserDB
)

var (
	flDebug      *bool
	flConfigFile *string

	logger *golog.GoLogger
)

var UserTable map[string][]net.Addr // must be a race condition.

// Note: please delete this func if you can find a built-in one.
func addrInSlice(a net.Addr, list []net.Addr) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func addrPositionInSlice(a net.Addr, list []net.Addr) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}

func init() {
	flDebug = flag.Bool("d", false, "debug switch")
	flConfigFile = flag.String("c", "./config.json", "config file")
}

func handleConnection(c net.Conn) {
	r, err := tunnel.NewServer(sc.Tunnel, sc.MasterKey, sc.EncryptMethod, logger)
	if err != nil {
		logger.Fatalf("Create tunnel failed: %s", err)
	}

	r.Accept(c)

	s := protocol.NewServer(nil)

	user := s.Accept(r)
	if UserTable[user] == nil {
		UserTable[user] = make([]net.Addr, 1)
		UserTable[user][0] = c.RemoteAddr()
	} else if !addrInSlice(c.RemoteAddr(), UserTable[user]) {
		// fixme: using configurable limits.
		if len(UserTable[user]) >= 3 {
			// fixme: adding debug info. Let user know why it failed.
			return
		}
		UserTable[user] = append(UserTable[user], c.RemoteAddr())
	}
	ui, ok := udb.GetUserInfo(user)
	if !ok {
		c.Close()
		return
	}

	r.SetPassword(ui.Password)

	_, addrPort, err := s.ParseUserHeader(r)
	if err != nil {
		logger.Errorf("ParseUserHeader failed: %v", err)
		c.Close()
		return
	}

	realServer, err := net.Dial("tcp", addrPort)

	if err != nil {
		logger.Errorf("err: %s", err)
		return
	}

	go s.Upstream(r, realServer)
	s.Downstream(r, realServer)

	realServer.Close()
	pos := addrPositionInSlice(c.RemoteAddr(), UserTable[user])
	UserTable[user] = append(UserTable[user][:pos], UserTable[user][pos+1:]...)
	c.Close()
}

func main() {
	var err error
	UserTable = make (map[string][]net.Addr)

	// FIXME: configurable logger file
	logger = golog.New(os.Stdout)

	flag.Parse()

	logger.SetDebug(*flDebug)

	// FIXME: configurable threads
	runtime.GOMAXPROCS(9)

	sc, err = parseConfigFile(*flConfigFile)
	if err != nil {
		logger.Fatalf("Parse config file err: %s", err)
	}

	udb, err = userdb.New(sc.UserDB)
	if err != nil {
		logger.Fatalln(err)
	}

	udb.SyncFromDB()

	lnTCP, err := net.Listen("tcp", fmt.Sprintf("%s:%d", sc.ListenAddr, sc.ListenPort))
	if err != nil {
		panic(err)
	}
	defer lnTCP.Close()

	for {
		conn, err := lnTCP.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}
