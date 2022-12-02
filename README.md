# `taskmaster`
Process control system in Go

> **Warning**
>
> This project is currently in progress

## About
`taskmaster` is a school project, part of my cursus at [19](https://campus19.be). The directives are to write a supervisor-like process manager in any language with a server running child processes based on a configuration file and a client which sends basic instruction (restart, reload, kill, etc) in any language.

## Design decisions
I decided to write the project in Go since it is a language I am currently learning and wish to become better at.
For the the client-server communcation, I went for Google's [gRPC](https://grpc.io/) through UNIX sockets since it fits my needs perfectly.

## Usage

### Server
```
Usage of taskmasterd
  -conf string
        config file path (default "/etc/taskmaster.yaml")
  -socket string
        path of the server's unix socket (default "/tmp/taskmaster.sock")
```
Example:
```
$ taskmasterd -conf ./prod.yaml -socket /tmp/taskmaster-prod.sock
```

### Client
```
Usage of taskmasterctl
  -socket string
        path of the server's unix socket (default "/tmp/taskmaster.sock")
```
Example:
```
$ taskmasterctl -socket /tmp/taskmaster-prod.sock
tm> service nginx status
SERVICE: nginx
STATUS: running
```

### Example configuration file
```
programs:
    nginx:
        cmd: "/usr/local/bin/nginx -c /etc/nginx/test.conf"
        numprocs: 1
        umask: 022
        workingdir: /tmp
        autostart: true
        autorestart: unexpected
        exitcodes:
            - 0
            - 2
        startretries: 3
        starttime: 5
        stopsignal: TERM
        stoptime: 10
        stdout: /tmp/nginx.stdout
        stderr: /tmp/nginx.stderr
        env:
        STARTED_BY: taskmaster
        ANSWER: 42
    vogsphere:
        cmd: "/usr/local/bin/vogsphere-worker --no-prefork"
        numprocs: 8
        umask: 077
        workingdir: /tmp
        autostart: true
        autorestart: unexpected
        exitcodes:
            - 0
        startretries: 3
        starttime: 5
        stopsignal: USR1
        stoptime: 10
        stdout: /tmp/vgsworker.stdout
        stderr: /tmp/vgsworker.stderr

```