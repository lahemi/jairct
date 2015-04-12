package main

import "os"

var homePath = os.Getenv("HOME")

var configPath = homePath + "/.config/jairct"
var dataPath = homePath + "/.local/share/jairct"

var dbFile = dataPath + "/db"

var defaultPort = "6667"

var initNetwork = "irc.freenode.net"
var initNick = "jairctwat"
