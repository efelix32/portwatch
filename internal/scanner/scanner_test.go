package scanner

import (
	"testing"

	psnet "github.com/shirou/gopsutil/v3/net"
)

func TestFromConnectionStats_MapsFields(t *testing.T) {
	stats := []psnet.ConnectionStat{
		{
			Type:   1, // tcp
			Laddr:  psnet.Addr{IP: "127.0.0.1", Port: 8080},
			Raddr:  psnet.Addr{IP: "0.0.0.0", Port: 0},
			Status: "LISTEN",
			Pid:    -1, // guaranteed not to resolve to a real process
		},
	}

	conns := FromConnectionStats(stats)

	if len(conns) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(conns))
	}
	got := conns[0]
	if got.Proto != "tcp" {
		t.Errorf("Proto = %q, want tcp", got.Proto)
	}
	if got.LocalAddr != "127.0.0.1:8080" {
		t.Errorf("LocalAddr = %q, want 127.0.0.1:8080", got.LocalAddr)
	}
	if got.Status != "LISTEN" {
		t.Errorf("Status = %q, want LISTEN", got.Status)
	}
	if got.ProcessName != "" {
		t.Errorf("ProcessName = %q, want empty for unresolvable pid", got.ProcessName)
	}
}

func TestFromConnectionStats_UnknownProtoType(t *testing.T) {
	stats := []psnet.ConnectionStat{{Type: 99, Pid: -1}}

	conns := FromConnectionStats(stats)

	if conns[0].Proto != "?" {
		t.Errorf("Proto = %q, want ?", conns[0].Proto)
	}
}

func TestSortByLocalPort_OrdersNumericallyNotLexically(t *testing.T) {
	// 99 < 100 numerically, but "100" < "99" as strings -- this case
	// catches a naive string-sort implementation.
	conns := []Connection{
		{LocalAddr: "127.0.0.1:9000"},
		{LocalAddr: "127.0.0.1:100"},
		{LocalAddr: "127.0.0.1:99"},
		{LocalAddr: "127.0.0.1:80"},
	}

	SortByLocalPort(conns)

	want := []string{"127.0.0.1:80", "127.0.0.1:99", "127.0.0.1:100", "127.0.0.1:9000"}
	for i, w := range want {
		if conns[i].LocalAddr != w {
			t.Errorf("position %d: got %q, want %q", i, conns[i].LocalAddr, w)
		}
	}
}

func TestSortByLocalPort_MalformedAddrSortsFirst(t *testing.T) {
	conns := []Connection{
		{LocalAddr: "127.0.0.1:80"},
		{LocalAddr: "not-an-address"},
	}

	SortByLocalPort(conns)

	if conns[0].LocalAddr != "not-an-address" {
		t.Errorf("expected malformed address to sort first, got %q", conns[0].LocalAddr)
	}
}
