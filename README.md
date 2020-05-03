# SMPP Gateway

<a href="https://goreportcard.com/report/github.com/xhit/go-simple-mail/v2"><img src="https://goreportcard.com/badge/github.com/xhit/go-simple-mail" alt="Go Report Card"></a>

## Description

SMPP clients testing tool.

Runs SMPP server, forwards incoming SMS with some debug information to SMTP.

Based on [SMPP 3.4 Library](https://github.com/ajankovic/smpp "SMPP 3.4 Library") 


## Goals
- Receiving SMS and forwarding them to an SMTP server for testing client applications using SMPP in simple scenarios.

## Features
- SMPP supports protocol version: **3.4**
- SMPP supports PDU commands: **bind_transceiver, submit_sm, unbind**
- SMPP supports encodings: **UCS-2**
- SMTP supports authentication methods: **PLAIN,  LOGIN, CRAM-MD5**
- SMTP supports encryption methods: **None,  SSL/TLS, STARTTLS**

## SMPP docs
[Short Message Peer to Peer
Protocol Specification v3.4](http://docs.nimta.com/SMPP_v3_4_Issue1_2.pdf)

[SMPP v3.4 Protocol Implementation
guide for GSM / UMTS](http://opensmpp.org/specs/smppv34_gsmumts_ig_v10.pdf)


## License
[MIT](LICENSE)
