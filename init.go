package brook

import (
	ant "github.com/txthinking/ant"
	"github.com/txthinking/socks5"
)

var Dial ant.Dialer = ant.DefaultDial

// EnableDebug
func EnableDebug() {
	socks5.Debug = true
}
