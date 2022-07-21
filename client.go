package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

type client struct {
    conn net.Conn
    name string
    group *group // for now, let's just have one chat group per user
    actions chan<- command
}

func (c *client) msg(msg string) {
    c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *client) err(err error) {
    c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) readActions() {
    for {
        msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")

        switch args[0] {
        case "/name":
            c.passCommand(NAME, args)
        case "/join":
            c.passCommand(JOIN, args)
        case "/groups":
            c.passCommand(GROUPS, args)
        case "/msg":
            c.passCommand(MSG, args)
        case "/abort":
            c.passCommand(ABORT, args)
        case "/help":
            c.printHelp()
        default:
            ok := c.msgIfInRoom(args)
            if !ok {
                c.err(fmt.Errorf("unknown command: %s", args[0]))
            }
        }
    }
}

func (c *client) msgIfInRoom(args []string) bool {
    if c.group != nil {
        args = append([]string{"msg"}, args...)
        c.passCommand(MSG, args)
    }
    return c.group != nil
}

func (c *client) printHelp() {
    c.msg("Available commands: /name, /join, /groups, /msg, /abort, /help")
}

func (c *client) passCommand(cmd commandUUID, args []string) {
    c.actions <- command {
        uuid: cmd,
        args: args,
        client: c,
    }
}
