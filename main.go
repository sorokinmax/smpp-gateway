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
var version = "v.1.0.0"

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
	flag.StringVar(&systemID, "systemid", "Sorokin's SMS Gateway", "descriptive server identification.")
	flag.Parse()

	sessConf := smpp.SessionConf{
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
			switch ctx.CommandID() {
			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					fail("Invalid PDU in context error: %+v", err)
				}
				resp := btrx.Response(systemID)
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					fail("Server can't respond to the Binding request: %+v", err)
				}

			case pdu.SubmitSmID:
				sm, err := ctx.SubmitSm()
				if err != nil {
					fail("Invalid PDU in context error: %+v", err)
				}
				msg := `
				<html>
					<head>
						<meta http-equiv="Content-Type" content="text/html; charset=\"utf-8\"">
					</head>
					<body>
						From: ` + sm.SourceAddr + `<br>
						To: ` + sm.DestinationAddr + `<br>
						Priority: ` + strconv.Itoa(sm.PriorityFlag) + `<br>
						SMS: ` + UCS2Decode(sm.ShortMessage) + `<br><br><br><hr>
						SMPP Gateway ` + version + `<br>
						<i>Пупсик, если тебе что-то не нравится - пиши:</i> <a href="https://github.com/sorokinmax" target="_blank">Maxim Sorokin</a>
					</body>
				</html>
				`
				log.Println(fmt.Sprintf("MSG: %s\n", msg))
				SMTPSend("SMPP gateway", msg)

				msgID++
				resp := sm.Response(fmt.Sprintf("msgID_%d", msgID))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					fail("Server can't respond to the submit_sm request: %+v", err)
				}
			case pdu.UnbindID:
				unb, err := ctx.Unbind()
				if err != nil {
					fail("Invalid PDU in context error: %+v", err)
				}
				resp := unb.Response()
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					fail("Server can't respond to the submit_sm request: %+v", err)
				}
				ctx.CloseSession()
			}
		}),
	}
	srv := smpp.NewServer(serverAddr, sessConf)

	log.Println(fmt.Sprintf("'%s' is listening on '%s'\n", systemID, serverAddr))
	err = srv.ListenAndServe()
	if err != nil {
		fail("Serving exited with error: %+v", err)
	}
	log.Println(fmt.Sprintf("Server closed\n"))
}

func fail(msg string, params ...interface{}) {
	log.Println(fmt.Sprintf(msg+"\n", params...))
	os.Exit(1)
}
