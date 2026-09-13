// Package scanner reads the local machine's active network connections.
package scanner

import (
	"fmt"
	"net"
	"sort"
	"strconv"

	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Connection is a single, display-ready network connection row.
type Connection struct {
	Proto       string
	LocalAddr   string
	RemoteAddr  string
	Status      string
	PID         int32
	ProcessName string
}

// Scan returns the current TCP/UDP connections on the machine, sorted by
// local port so the table doesn't jump around between refreshes.
func Scan() ([]Connection, error) {
	stats, err := psnet.Connections("inet")
	if err != nil {
		return nil, fmt.Errorf("scanner: reading connections: %w", err)
	}
	return FromConnectionStats(stats), nil
}

// FromConnectionStats converts gopsutil's raw connection stats into the
// package's own Connection type. Split out from Scan so it can be unit
// tested without touching the real network stack.
func FromConnectionStats(stats []psnet.ConnectionStat) []Connection {
	conns := make([]Connection, 0, len(stats))
	for _, s := range stats {
		conns = append(conns, Connection{
			Proto:       protoName(s.Type),
			LocalAddr:   fmt.Sprintf("%s:%d", s.Laddr.IP, s.Laddr.Port),
			RemoteAddr:  fmt.Sprintf("%s:%d", s.Raddr.IP, s.Raddr.Port),
			Status:      s.Status,
			PID:         s.Pid,
			ProcessName: processName(s.Pid),
		})
	}
	SortByLocalPort(conns)
	return conns
}

// SortByLocalPort orders connections by their local port, ascending
// (numerically -- a plain string sort would put "9000" before "80").
func SortByLocalPort(conns []Connection) {
	sort.Slice(conns, func(i, j int) bool {
		return localPort(conns[i].LocalAddr) < localPort(conns[j].LocalAddr)
	})
}

// localPort extracts the numeric port from a "host:port" address, treating
// anything unparsable as port 0 so malformed addresses sort first rather
// than panicking or silently misordering the rest of the table.
func localPort(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}

func protoName(t uint32) string {
	switch t {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	default:
		return "?"
	}
}

// processName looks up a process's executable name by PID, returning an
// empty string (never an error) when the process can't be resolved --
// it may have exited between the scan and this lookup.
func processName(pid int32) string {
	if pid <= 0 {
		return ""
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return ""
	}
	name, err := p.Name()
	if err != nil {
		return ""
	}
	return name
}
