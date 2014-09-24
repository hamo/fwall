package main

import (
	"flag"
	"fmt"
	"net"
	"os"

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

func init() {
	flDebug = flag.Bool("d", false, "debug switch")
	flConfigFile = flag.String("c", "./config.json", "config file")
}

func handleConnection(c net.Conn) {
	r, err := tunnel.NewRawSocketServer(sc.MasterKey, sc.EncryptMethod, logger)
	if err != nil {
		logger.Fatalf("Create tunnel failed: %s", err)
	}

	r.Accept(c)

	s := protocol.NewServer(nil)

	user := s.Accept(r)
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
		fmt.Printf("err: %s", err)
		return
	}

	go s.Upstream(r, realServer)
	s.Downstream(r, realServer)

	realServer.Close()
	c.Close()
}

func main() {
	var err error

	// FIXME: configurable logger file
	logger = golog.New(os.Stdout)

	flag.Parse()

	logger.SetDebug(*flDebug)

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
