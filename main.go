package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
)

var (
	serverAddr string
	systemID   string
	msgID      int
)
var version = "v1.5.0"

var cfg Config

func main() {
	readConfigFile(&cfg)

	log.SetFlags(log.LstdFlags)
	f, err := os.OpenFile("smpp-gateway.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	multi := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multi)

	serverHost := cfg.SMPP.Host + ":" + strconv.Itoa(cfg.SMPP.Port)
	flag.StringVar(&serverAddr, "addr", serverHost, "server will listen on this address.")
	flag.StringVar(&systemID, "systemid", "SMS Gateway", "descriptive server identification.")
	flag.Parse()

	go CacheAutoUpdater("./mapping.txt")

	sessConf := smpp.SessionConf{
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {

			switch ctx.CommandID() {
			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					Logger("Invalid PDU in context error: %+v", err)
				}
				Logger("Incoming connection from %s with ID: %s", ctx.RemoteAddr(), btrx.SystemID)
				resp := btrx.Response(systemID)
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the Binding request: %+v", err)
				}

			case pdu.SubmitSmID:
				sm, err := ctx.SubmitSm()
				if err != nil {
					Logger("Invalid PDU in context error: %+v", err)
				}

				SendSMS(sm, ctx)

				resp := sm.Response(fmt.Sprintf("msgID_%d", msgID))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the submit_sm request: %+v", err)
				}
				msgID++

			case pdu.UnbindID:
				unb, err := ctx.Unbind()
				if err != nil {
					Logger("Invalid PDU in context error: %+v", err)
				}
				resp := unb.Response()
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the submit_sm request: %+v", err)
				}
				ctx.CloseSession()
			}
		}),
	}
	srv := smpp.NewServer(serverAddr, sessConf)

	Logger("'%s' is listening on '%s'", systemID, serverAddr)
	err = srv.ListenAndServe()
	if err != nil {
		Logger("Serving exited with error: %+v", err)
	}
	Logger("Server closed")
}

// Logger - logging wrapper
func Logger(msg string, params ...interface{}) {
	log.Printf(msg, params...)
}
