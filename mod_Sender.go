package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/sorokinmax/smpp"
	"github.com/sorokinmax/smpp/pdu"
)

func SendEmail(sm *pdu.SubmitSm, ctx *smpp.Context, message string) (responce string, status pdu.Status) {
	msgEmail := `
			<body>
				SessionId: ` + ctx.SessionID() + `<br>
				From: ` + sm.SourceAddr + `<br>
				To: ` + sm.DestinationAddr + `<br>
				Priority: ` + strconv.Itoa(sm.PriorityFlag) + `<br>
				RemoteAddress: ` + ctx.RemoteAddr() + `<br><br>
				SMS: ` + message + `<br><br><br><hr>
				<a href="https://github.com/sorokinmax" target="_blank">SMPP Gateway ` + version + `</a><br>
			</body>
			`

	err := SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Encr, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From, cache.Mapping[sm.DestinationAddr], "SMPP gateway", msgEmail, "")
	if err != nil {
		log.Printf("msgID %d: sending error", msgID)
		responce = fmt.Sprintf("msgID_%d: sending error /n %s", msgID, err.Error())
		status = pdu.StatusDeliveryFailure
	} else {
		log.Printf("msgID %d: successfully sent", msgID)
		responce = fmt.Sprintf("msgID_%d: successfully sent", msgID)
		status = pdu.StatusOK
	}
	return responce, status
}

func SendTelegram(sm *pdu.SubmitSm, ctx *smpp.Context, message string) (responce string, status pdu.Status) {
	msgTg := `	SessionId: ` + ctx.SessionID() + `
			From: ` + sm.SourceAddr + `
			To: ` + sm.DestinationAddr + `
			Priority: ` + strconv.Itoa(sm.PriorityFlag) + `
			RemoteAddress: ` + ctx.RemoteAddr() + `

			SMS: ` + message + `

			==================
			SMPP Gateway ` + version + `
			`

	// Convert string to int
	dst, err := strconv.Atoi(cache.Mapping[sm.DestinationAddr])
	if err != nil {
		log.Printf("msgID_%d: not sent.", msgID)
		resp := sm.Response(fmt.Sprintf("msgID_%d: not sent. /n %s", msgID, err.Error()))
		if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
			log.Printf("Server can't respond to the submit_sm request: %+v", err)
		}
	}
	// Sending
	_, err = tgSendMessage(cfg.Telegram.BotToken, msgTg, dst)
	if err != nil {
		log.Printf("msgID_%d: sending error", msgID)
		responce = fmt.Sprintf("msgID_%d: sending error /n %s", msgID, err.Error())
		status = pdu.StatusDeliveryFailure
	} else {
		log.Printf("msgID_%d: successfully sent", msgID)
		responce = fmt.Sprintf("msgID_%d: successfully sent", msgID)
		status = pdu.StatusOK
	}

	return responce, status
}
