# prockeeper

A process manager app written in go.

<img src="https://raw.githubusercontent.com/jiajiawang/prockeeper/master/prockeeper-preview.png" width=640>

## Installation

`go install github.com/jiajiawang/prockeeper@latest`

## Usage

`prockeeper --help`

```
usage: prockeeper [options]

  --help          Show this help
  --debug         Write application log to ./prockeeper.log
  -c path_to_yml  Specify the path of yaml file (default: './prockeeper.yml')

Service Options:
  [name]    Specify the name of the service
  [command] Specify the exec command
  [dir]     Specify the working directory

Example yaml:
  services:
    - name: "rails server"
      command: "rails s"
    - name: "node server"
      command: "npm start"
      dir: "./client"

Keyboard commands

j      - Select previous item
k      - Select next item
Enter  - Start/stop selected service
u      - Start all services
d      - Stop all services

?      - Show/hide help menu
.      - Show/hide application log
Ctrl-C - Exit app
```
