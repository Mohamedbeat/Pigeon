package main

import "errors"

const (
	ConnectCMD     = "CONNECT"
	SubscribeCMD   = "SUBSCRIBE"
	UnsubscribeCMD = "UNSUBSCRIBE"
	PublishCMD     = "PUBLISH"
	DisconnectCMD  = "DISCONNECT"
	ListCMD        = "LIST"
	PingCMD        = "PING"
)

var UnknownCommandErr = errors.New("unknown command")
var MalformedCommandErr = errors.New("Malformed command")
var PeerNotFoundErr = errors.New("Peer not found")
var TopicNotFoundErr = errors.New("Topic not found")
