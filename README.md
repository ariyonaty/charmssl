## CharmSSL

A first take at a more user-friendly x509 Certificate Viewer using Charm's [BubbleTea](https://github.com/charmbracelet/bubbletea)

Note to reader: this is more of a small project to get some hands-on experience with Go and one of its popular libraries. Issues, suggestions, or contributions are more than welcome.

### Usage

#### Local Certificate

```
$ go run main.go -file github.pem
```

#### Remote Certificate
```
$ go run main.go -domain google.com
```
