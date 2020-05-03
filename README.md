# prockeeper

A process manager app written in go.

<img src="https://raw.githubusercontent.com/jiajiawang/prockeeper/master/prockeeper-preview.png" width=640>

## Installation

`go get github.com/jiajiawang/prockeeper`

## Usage

Add `prockeeper.yml` to project root and run `prockeeper`.
Or use `prockeeper -c path_to_yml` to specify the config file.

Example yaml:
```
services:
  - name: "rails server"
  command: "rails s"
  - name: "node server"
  command: "npm start"
```

## Keyboard commands

```
j      - Select previous item
k      - Select next item
Enter  - Start/stop selected service
u      - Start all services
d      - Stop all services

?      - Show/hide help menu
.      - Show/hide debugger
Ctrl-C - Exit app
```
