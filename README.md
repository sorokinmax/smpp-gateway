# SMPP Gateway


## Description

SMPP clients testing tool.

Runs SMPP server, forwards incoming SMS with some debug information to SMTP or Telegram.

Incoming SMS are redirected to emails or Telegram, mapped with phone numbers in mapping.txt

Based on [Modified  SMPP 3.4 Library](https://github.com/sorokinmax/smpp)


## Goals
- Receiving SMS and forwarding them to an SMTP server or Telegram for testing client applications using SMPP in simple scenarios.

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
