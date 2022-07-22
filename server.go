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
        groups: make(map[string]*group),
        activeGroup: nil,
        actions: s.commands,
    }
}

func (s *server) run() {
    for command := range s.commands {
        switch command.uuid {
        case NAME:
            if s.checkCommand(command) {
                s.name(command.client, command.args[1])
            }
        case JOIN:
            if s.checkCommand(command) {
                s.join(command.client, command.args[1])
            }
        case GROUPS:
            s.allGroups(command.client)
        case MSG:
            s.msg(command.client, command.args)
        case ABORT:
            s.abort(command.client)
        case SWITCH:
            if s.checkCommand(command) {
                s.switchGroup(command.client, command.args[1])
            }
        case LEAVE:
            s.leaveGroup(command.client, command.args)
        }
    }
}

func (s *server) checkCommand(c command) bool {
    if len(c.args) < 2 {
        c.client.msg(fmt.Sprintf("Improper %[1]s usage; use /%[1]s <%[1]s>.", c.args[0][1:]))
    }
    return len(c.args) >= 2
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

    c.groups[g.title] = g
    c.activeGroup = g

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
	c.activeGroup.broadcast(c, c.name + ": " + msg)
}

func (s *server) abort(c *client) {
    log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
    s.kickClient(c)
    c.msg("We hope you enjoyed your stay. Bye!")
    c.conn.Close()
}

func (s *server) kickClient(c *client) {
    if c.activeGroup != nil {
		temp := s.groups[c.activeGroup.title]
		delete(s.groups[c.activeGroup.title].members, c.conn.RemoteAddr())
		temp.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}

func (s *server) switchGroup(c *client, groupName string) {
    if val, ok := c.groups[groupName]; ok {
        // client is in the group => able to set activeGroup to group
        c.activeGroup = val
        c.msg(fmt.Sprintf("You have switched to the group: %s", groupName))
    } else {
        c.msg(fmt.Sprintf("You are not in the group: %s. Try the /join command.", groupName))
    }
}

func (s *server) leaveGroup(c *client, args []string) {
    var groupName string
    if len(c.groups) == 0 {
        c.msg("You are not a member of any groups.")
        return
    }
    if len(args) < 2 {
        // leave current group
        groupName = c.activeGroup.title
        s.kickClient(c)
        c.activeGroup = nil
        delete(c.groups, groupName)
    } else {
        // leave the group specified by args[1]
        groupName = args[1]
        s.kickClient(c)
        delete(c.groups, groupName)
    }
    c.msg(fmt.Sprintf("You have left the group: %s", groupName))
}

