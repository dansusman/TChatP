package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
    groups map[string]*group
    commands chan command
}

func newServer() *server {
    return &server {
        groups: make(map[string]*group),
        commands: make(chan command),
    }
}

func (s *server) newClient(conn net.Conn) *client {
    return &client {
        name: "anon",
        conn: conn,
        group: nil,
        actions: s.commands,
    }
}

func (s *server) run() {
    for command := range s.commands {
        switch command.uuid {
        case NAME:
            s.name(command.client, command.args[1])
        case JOIN:
            s.join(command.client, command.args[1])
        case GROUPS:
            s.allGroups(command.client)
        case MSG:
            s.msg(command.client, command.args)
        case ABORT:
            s.abort(command.client)
        }
    }
}

func (s *server) name(c *client, name string) {
    c.name = name
    c.msg(fmt.Sprintf("nice name, %s", c.name))
}

func (s *server) join(c *client, groupName string) {
    g, ok := s.groups[groupName]

    if !ok {
        g = &group {
            title: groupName,
            members: make(map[net.Addr]*client),
        }
        s.groups[groupName] = g
    }

    otherMembers := g.memberNames()
    g.members[c.conn.RemoteAddr()] = c

    g.broadcast(c, fmt.Sprintf("%s has joined the group!", c.name))

    c.group = g

    c.msg(fmt.Sprintf("%s welcomes you, %s! enjoy your stay.", groupName, c.name))

    if len(otherMembers) != 0 {
        c.msg(fmt.Sprintf("other members of this group: %s", strings.Join(otherMembers, ", ")))
    }
}

func (s *server) allGroups(c *client) {
    var groupNames []string

    for name := range s.groups {
        groupNames = append(groupNames, name)
    }
    c.msg(fmt.Sprintf("available chatgroups: %s", strings.Join(groupNames, ", ")))
}

func (s *server) msg(c *client, args []string) {
    if len(args) < 2 {
		c.msg("message is required, usage: /msg <some text here>")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.group.broadcast(c, c.name + ": " + msg)
}

func (s *server) abort(c *client) {
    log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
    s.kickClient(c)
    c.msg("We hope you enjoyed your stay. Bye!")
    c.conn.Close()
}

func (s *server) kickClient(c *client) {
    if c.group != nil {
		temp := s.groups[c.group.title]
		delete(s.groups[c.group.title].members, c.conn.RemoteAddr())
		temp.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}

