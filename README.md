# aws-tools
A collection of shortcuts to make common AWS operations easier.

## aws-ssh

### Install
`go get github.com/oremj/aws-tools/aws-ssh`

### Usage
`aws-ssh [-c COMMAND] instance-id [ssh-args]`

## aws-instancelist

### Install
`go get github.com/oremj/aws-tools/aws-instancelist`

### Usage
`aws-instancelist [-f filter]...`

### Example
`aws-instancelist -f "tag:App=myapp" -f "tag:Type=web,admin"`
