package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"

	"github.com/efelix32/portwatch/internal/scanner"
)

func TestRowsFromConnections(t *testing.T) {
	conns := []scanner.Connection{
		{Proto: "tcp", LocalAddr: "127.0.0.1:8080", RemoteAddr: "0.0.0.0:0", Status: "LISTEN", PID: 1234, ProcessName: "node"},
	}

	rows := rowsFromConnections(conns)

	want := table.Row{"tcp", "127.0.0.1:8080", "0.0.0.0:0", "LISTEN", "1234", "node"}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	got := rows[0]
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("column %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRowsFromConnections_Empty(t *testing.T) {
	rows := rowsFromConnections(nil)
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}
