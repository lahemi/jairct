package main

import (
	"regexp"
	"strings"
)

var lineregex = regexp.MustCompile("^:[^ ]+?!([^ ]+? ){3}:.+")

type MsgLine struct {
	Nick, Cmd, Target, Msg string
}

func SendToCan(can, line string) {
	writechan <- "PRIVMSG " + can + " :" + line
}

func SplitMsgLine(l string) MsgLine {
	spl := strings.SplitN(l, " ", 4)
	return MsgLine{
		Nick:   spl[0][1:strings.Index(l, "!")],
		Cmd:    spl[1],
		Target: spl[2],
		Msg:    spl[3][1:],
	}
}

func HandleOut(s string) {
	if lineregex.MatchString(s) {
		ml := SplitMsgLine(s)
		sep := " | "
		stdout(ml.Nick + sep + ml.Target + sep + ml.Msg)
	} else {
		stdout(s)
	}
}
