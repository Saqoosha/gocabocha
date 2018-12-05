package cabocha

import (
	"testing"
)

func TestChunking(t *testing.T) {
	p := "太郎 は花子が、持っている .本を\t次郎に渡した！"
	expected := []string{"太郎 は", "花子が、", "持っている ", ".本を\t", "次郎に", "渡した！"}
	for i, c := range Chunks(p) {
		if c != expected[i] {
			t.Fatalf("not expected! %v %v", i, c)
		}
	}
}
