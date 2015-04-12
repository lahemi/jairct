package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

var writechan = make(chan string, 512)

func connectServer(server, port string) (conn net.Conn) {
	conn, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		die(err)
	}
	return
}

func readServer(r *bufio.Reader) {
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			stderr("read error")
			break
		}
		if str[:4] == "PING" {
			writechan <- "PONG" + str[4:len(str)-2]
		} else {
			line := str[:len(str)-2]
			//HandleOut(line)
			stdout(line)
			if lineregex.MatchString(line) {
				msg := SplitMsgLine(line)
				go saveDB(msg)
			}
		}
	}
}

func writeServer(w *bufio.Writer) {
	for {
		str := <-writechan
		if _, err := w.WriteString(str + "\r\n"); err != nil {
			stderr("write error")
			break
		}
		// TODO: proper check for PRIVMSG
		if strings.Contains(str, "PRIVMSG") {
			go saveDB(SplitMsgLine(str))
		}
		stdout(str)
		w.Flush()
	}
}

// Is setup fails, application exits with error message and $? of 1.
func init() {
	setupEnvironment()
}

func main() {
	conn := connectServer(initNetwork, defaultPort)

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	in := bufio.NewReader(os.Stdin)

	go readServer(r)
	go writeServer(w)

	writechan <- "USER " + initNick + " * * :" + initNick
	writechan <- "NICK " + initNick

	for {
		input, err := in.ReadString('\n')
		if err != nil {
			stderr("error input")
			break
		}
		inp := input[:len(input)-1]

		var quitSig bool
		if inp == "QUIT" {
			quitSig = true
		}

		if len(inp) != 0 {
			if quitSig {
				DB.Close()
			}
			writechan <- inp
			if quitSig {
				err = conn.Close()
				if err != nil {
					die("Connection failed to close: " + err.Error())
				}
				// Need a little wait here, or some check to ensure all the other closing is done..?
				break
			}
		}
	}
}
