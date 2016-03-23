// +build windows

package pb

import (
	"gopkg.in/andrewgoktepe/etcd.v2/Godeps/_workspace/src/github.com/olekukonko/ts"
)

func bold(str string) string {
	return str
}

func terminalWidth() (int, error) {
	size, err := ts.GetSize()
	return size.Col(), err
}
