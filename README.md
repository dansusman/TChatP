# Description

This project is a simple TCP chat server that allows users to communicate
back and forth with other members of their chat group. My main motivation
is to learn more sophisticated Go.


# Supported Commands

1. `/name <name>`: give yourself a name; if not used, a user is considered anonymous.
2. `/join <groupName>`: join a chat group; if the group doesn't exist, a new
one with the specified name will be created.
    - Note: each user can be a member of one group at a time.
3. `/groups`: show the list of all available chat groups.
4. `/msg <msg>`: broadcast a message to all members of the chat group.
5. `/abort`: disconnect from the chat server.
