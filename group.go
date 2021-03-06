package main

import "net"

type group struct {
	title   string
	members map[net.Addr]*client
}

func (g *group) broadcast(sender *client, msg string) {

	for addr, cl := range g.members {
		if sender.conn.RemoteAddr() != addr {
			cl.msg(msg)
		}
	}
}

func (g *group) memberNames() []string {
	var names []string
	for _, cl := range g.members {
		names = append(names, cl.name)
	}
	return names
}
