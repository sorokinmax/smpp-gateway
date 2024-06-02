package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"strconv"

	"github.com/sorokinmax/smpp"
	"github.com/sorokinmax/smpp/pdu"
)

const version = "v1.7.1"

/*
type AppRegistry struct {
	SessionId string
	SystemId  string
	Password  string
}*/

var (
	serverAddr string
	systemID   string
	msgID      int
	//appData    map[string]AppRegistry
)

var cfg Config

func init() {
	readConfigFile(&cfg)
	Messages = make(map[byte]map[byte]Message)
}

func main() {
	log.SetFlags(log.LstdFlags)
	f, err := os.OpenFile("smpp-gateway.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	multi := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multi)

	//appData = make(map[string]AppRegistry)

	serverHost := cfg.SMPP.Host + ":" + strconv.Itoa(cfg.SMPP.Port)
	flag.StringVar(&serverAddr, "addr", serverHost, "server will listen on this address.")
	flag.StringVar(&systemID, "systemid", "SMS Gateway", "descriptive server identification.")
	flag.Parse()

	go CacheAutoUpdater("./mapping.txt")

	sessConf := smpp.SessionConf{
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {

			switch ctx.CommandID() {
			case pdu.BindTransmitterID:
				btx, err := ctx.BindTx()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				log.Printf("Incoming connection from %s with ID: %s", ctx.RemoteAddr(), btx.SystemID)

				/*appData[ctx.SessionID()] = AppRegistry{
					SessionId: ctx.SessionID(),
					SystemId:  btx.SystemID,
					Password:  btx.Password,
				}*/

				resp := btx.Response(systemID)
				responseStatus := pdu.StatusInvPaswd
				if (btx.Password == cfg.SMPP.Password && btx.SystemID == cfg.SMPP.User) || cfg.SMPP.User == "" {
					responseStatus = pdu.StatusOK
				}
				if err := ctx.Respond(resp, responseStatus); err != nil {
					log.Printf("Server can't respond to the Binding request: %+v", err)
				}
			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				log.Printf("Incoming connection from %s with ID: %s", ctx.RemoteAddr(), btrx.SystemID)

				/*appData[ctx.SessionID()] = AppRegistry{
					SessionId: ctx.SessionID(),
					SystemId:  btrx.SystemID,
					Password:  btrx.Password,
				}*/

				resp := btrx.Response(systemID)
				responseStatus := pdu.StatusInvPaswd
				if (btrx.Password == cfg.SMPP.Password && btrx.SystemID == cfg.SMPP.User) || cfg.SMPP.User == "" {
					responseStatus = pdu.StatusOK
				}
				if err := ctx.Respond(resp, responseStatus); err != nil {
					log.Printf("Server can't respond to the Binding request: %+v", err)
				}

			case pdu.SubmitSmID:
				var response string
				var status pdu.Status

				sm, err := ctx.SubmitSm()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}

				message := GetMessage(sm.ShortMessage)
				switch sm.DataCoding {
				case 8:
					message = UCS2Decode(message)
				}

				if len(message) > 0 {
					log.Printf("Incoming SMS\n\tFrom: %s\n\tTo:%s\n\tPriority: %s\n\tRemoteAddress: %s\n\tSMS: %s", sm.SourceAddr, sm.DestinationAddr, strconv.Itoa(sm.PriorityFlag), ctx.RemoteAddr(), message)
					if len(cache.Mapping[sm.DestinationAddr]) > 0 {
						log.Printf("Send SMS to: %s", cache.Mapping[sm.DestinationAddr])

						// Check dst type
						_, err := mail.ParseAddress(cache.Mapping[sm.DestinationAddr])
						if err == nil {
							response, status = SendEmail(sm, ctx, message)
						} else {
							response, status = SendTelegram(sm, ctx, message)
						}
					} else {
						log.Printf("msgID %d: address not matched.", msgID)
						response = fmt.Sprintf("msgID_%d: address not matched.", msgID)
						status = pdu.StatusDeliveryFailure
					}
				}

				resp := sm.Response(response)
				if err := ctx.Respond(resp, status); err != nil {
					log.Printf("Server can't respond to the submit_sm request: %+v", err)
				}
				msgID++

			case pdu.EnquireLinkID:
				el, err := ctx.EnquireLink()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}

				resp := el.Response()
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					log.Printf("Server can't respond to the enquire_link request: %+v", err)
				}

			case pdu.UnbindID:
				unb, err := ctx.Unbind()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				//delete(appData, ctx.SessionID())
				resp := unb.Response()
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					log.Printf("Server can't respond to the submit_sm request: %+v", err)
				}
				ctx.CloseSession()
			}
		}),
	}
	srv := smpp.NewServer(serverAddr, sessConf)

	log.Printf("'%s' is listening on '%s'", systemID, serverAddr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("Serving exited with error: %+v", err)
	}
	log.Print("Server closed")
}
