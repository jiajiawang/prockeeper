# prockeeper

A process manager app written in go.

## Installation

`go get github.com/jiajiawang/prockeeper`

## Usage

Add `prockeeper.yml` to project root and run `prockeeper`

Example yaml:
```
services:
  - name: "date"
  command: "while true; do date; sleep 1; done"
  - name: "time in milliseconds"
  command: "while true; do ruby -e 'puts Time.now.to_f'; sleep 0.1; done"
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
