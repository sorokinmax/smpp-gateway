package main

import (
	"fmt"
	"net/mail"
	"strconv"

	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
)

func SendSMS(sm *pdu.SubmitSm, ctx *smpp.Context) {
	msgEmail := `
	<body>
		From: ` + sm.SourceAddr + `<br>
		To: ` + sm.DestinationAddr + `<br>
		Priority: ` + strconv.Itoa(sm.PriorityFlag) + `<br>
		RemoteAddress: ` + ctx.RemoteAddr() + `<br><br>
		SMS: ` + UCS2Decode(sm.ShortMessage) + `<br><br><br><hr>
		<a href="https://github.com/sorokinmax" target="_blank">SMPP Gateway ` + version + `</a><br>
	</body>
`

	msgTg := `From: ` + sm.SourceAddr + `
To: ` + sm.DestinationAddr + `
Priority: ` + strconv.Itoa(sm.PriorityFlag) + `
RemoteAddress: ` + ctx.RemoteAddr() + `

SMS: ` + UCS2Decode(sm.ShortMessage) + `

==================
SMPP Gateway ` + version + `
`

	Logger("Incoming SMS\n\tFrom: %s\n\tTo:%s\n\tPriority: %s\n\tRemoteAddress: %s\n\tSMS: %s", sm.SourceAddr, sm.DestinationAddr, strconv.Itoa(sm.PriorityFlag), ctx.RemoteAddr(), UCS2Decode(sm.ShortMessage))

	if len(cache.Mapping[sm.DestinationAddr]) > 0 {
		Logger("Send SMS to: %s", cache.Mapping[sm.DestinationAddr])

		// Check dst type
		_, err := mail.ParseAddress(cache.Mapping[sm.DestinationAddr])
		if err == nil {
			// Send Email
			err := SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Encr, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From, cache.Mapping[sm.DestinationAddr], "SMPP gateway", msgEmail, "")
			if err != nil {
				Logger("msgID_%d: not sent.", msgID)
				resp := sm.Response(fmt.Sprintf("msgID_%d: not sent. /n %s", msgID, err.Error()))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the submit_sm request: %+v", err)
				}
			}
		} else {
			// Send Telegram message

			// Convert string to int
			dst, err := strconv.Atoi(cache.Mapping[sm.DestinationAddr])
			if err != nil {
				Logger("msgID_%d: not sent.", msgID)
				resp := sm.Response(fmt.Sprintf("msgID_%d: not sent. /n %s", msgID, err.Error()))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the submit_sm request: %+v", err)
				}
			}
			// Sending
			_, err = tgSendMessage(cfg.Telegram.BotToken, msgTg, dst)
			if err != nil {
				Logger("msgID_%d: not sent.", msgID)
				resp := sm.Response(fmt.Sprintf("msgID_%d: not sent. /n %s", msgID, err.Error()))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					Logger("Server can't respond to the submit_sm request: %+v", err)
				}
			}
		}
	} else {
		Logger("msgID_%d: address not matched.", msgID)
		resp := sm.Response(fmt.Sprintf("msgID_%d: address not matched.", msgID))
		if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
			Logger("Server can't respond to the submit_sm request: %+v", err)
		}
	}
}
