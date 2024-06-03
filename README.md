# SMPP Gateway


## Description

SMPP-clients testing tool.

Runs SMPP server, forwards incoming SMS with some debug information to SMTP or Telegram.

Incoming SMS are redirected to e-mail or Telegram, and are matched with recipients in mapping.txt

Based on [Modified  SMPP 3.4 Library](https://github.com/sorokinmax/smpp)


## Goals
- Receiving SMS and forwarding them to an SMTP server or Telegram for testing client applications using SMPP in simple scenarios.

## Features
**SMPP**
- Protocol version: **3.4**
- PDU commands: **bind_transceiver, bind_transmitter, submit_sm, enquire_link, unbind**
- Encodings: **UCS-2**
- Long messages: User data headers(UDH) method only.
**SMTP**
- Authentication methods: **PLAIN,  LOGIN, CRAM-MD5**
- Encryption methods: **None,  SSL/TLS, STARTTLS**
**Telegram**
  - Sending SMS via your bot 

## SMPP docs
[Short Message Peer to Peer
Protocol Specification v3.4](http://docs.nimta.com/SMPP_v3_4_Issue1_2.pdf)

[SMPP v3.4 Protocol Implementation
guide for GSM / UMTS](http://opensmpp.org/specs/smppv34_gsmumts_ig_v10.pdf)


## License
[MIT](LICENSE)
