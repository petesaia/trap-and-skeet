This is an experiment and still in development. 

```bash
$ trap-and-skeet

NAME:
   Trap and Skeet - Use to create containers that are meant to be abused.

USAGE:
   trap-and-skeet [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     fire, create     This will fire off an ephemeral container.
     destroy, remove  This will destroy all containers were created.
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

```bash
$ trap-and-skeet fire --help

NAME:
   trap-and-skeet fire - This will fire off an ephemeral container.

USAGE:
   trap-and-skeet fire [command options] [arguments...]

OPTIONS:
   --docker-image value          (default: "krlmlr/debian-ssh")
   --ssh-user value              (default: "root")
   --ssh-password value
   --ssh-private-key-file value
   --ssh-port value              (default: "2020")
```
