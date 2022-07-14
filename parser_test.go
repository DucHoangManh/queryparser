package queryparser

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type Name struct {
	First string
	Last  string
}

func (n *Name) QueryParse(in string) error {
	parts := strings.Split(in, "_")
	if len(parts) != 2 {
		return fmt.Errorf("invalid input")
	}
	n.First = parts[0]
	n.Last = parts[1]
	return nil
}

type Dest struct {
	Name Name
	Mot  float64
	Ba   int64
	IDs  []int
}

func TestQueryParser(t *testing.T) {
	req, err := http.NewRequest("GET", "/root?name=first_last&mot=1.123&ba=3&ids=1,4,6,9", nil)
	if err != nil {
		t.Error(err.Error())
	}
	var d Dest
	err = Parse(req, &d)
	if err != nil {
		t.Error(err.Error())
	}
}
