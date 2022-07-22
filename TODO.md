- [x] define types needed: commands, groups, client, server
- [x] add a main.go that spins up server and does other stuff
- [x] build out client to send messages to server to handle
    - utilize goroutines for blocking functionality
    - send commands via channel, server can process everything in the channel later
- [x] build out server to handle messages from clients and tell room to broadcast to members,
change user names, join groups, etc.
- [x] add functionality to room to handle message broadcasting

## Future Plans

### Should
- [x] make a /help command that prints what commands are available
- [ ] there's a bug when a member leaves and client tries to run a command
    - something to do with newline?

### Probably Will
- [x] list members of a group when you join? ("you can talk to x, y, z")
- [x] use /msg by default when you're in a room (so clients don't have to type /msg every time)
- [x] clients can join multiple groups
- [x] when clients can be in more than one group, make a /switch command to talk in another group
- [ ] make a /leave command to leave a group

### Could
- [ ] format text messages
- [ ] scrape for bad language in Family Friendly Mode
- [ ] when server closes connection, save to a db for later spin up
    - i.e. save across sessions
- [ ] autocomplete group names....!!!!!???
- [ ] i'm sure there are more commands that could be useful and fun to implement

### Probably Won't
- [ ] build a frontend if you want lol
