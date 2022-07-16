package queryparser

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type Age uint

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
	Headline    string `query:"title"`
	Age         Age
	Category    int64
	Length      float64
	IsPublished bool `query:"is_published"`
	Sub         *Name
	Names       []Name
}

func TestQueryParser(t *testing.T) {
	tests := []struct {
		url      string
		expected Dest
	}{
		{
			url:      "/root?title=golang",
			expected: Dest{Headline: "golang"},
		},
		{
			url:      "/root?age=10",
			expected: Dest{Age: 10},
		},
		{
			url:      "/root?category=4566",
			expected: Dest{Category: 4566},
		},
		{
			url:      "/root?length=145.6756",
			expected: Dest{Length: 145.6756},
		},
		{
			url:      "/root?is_published=true",
			expected: Dest{IsPublished: true},
		},
		{
			url:      "/root?sub=duc_hoang",
			expected: Dest{Sub: &Name{First: "duc", Last: "hoang"}},
		},
		{
			url: "/root?names=duc_hoang,toby_han,minh_le",
			expected: Dest{Names: []Name{
				{
					First: "duc",
					Last:  "hoang",
				},
				{
					First: "toby",
					Last:  "han",
				},
				{
					First: "minh",
					Last:  "le",
				},
			},
			},
		},
	}
	for _, test := range tests {
		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Error(err.Error())
		}
		var d Dest
		err = Parse(req, &d)
		if err != nil {
			t.Error(err.Error())
		}
		destEqual(t, test.expected, d)
	}
}

func destEqual(t *testing.T, expect, actual Dest) {
	assertEqual(t, expect.Headline, actual.Headline)
	assertEqual(t, expect.Age, actual.Age)
	assertEqual(t, expect.Category, actual.Category)
	assertEqual(t, expect.Length, actual.Length)
	assertEqual(t, expect.IsPublished, actual.IsPublished)
	asserEqualWithFunc(t, actual.Sub, expect.Sub, func(a, b *Name) bool {
		if a == nil && b != nil {
			return false
		}
		if a != nil && b == nil {
			return false
		}
		if a == nil && b == nil {
			return true
		}
		return nameEqualFunc(*a, *b)
	})
	assertEqual(t, len(expect.Names), len(actual.Names))
	for i := range expect.Names {
		asserEqualWithFunc(t, expect.Names[i], actual.Names[i], nameEqualFunc)
	}
}

func nameEqualFunc(a, b Name) bool {
	if a.First != b.First {
		return false
	}
	if a.Last != b.Last {
		return false
	}
	return true
}

func assertEqual[T comparable](t *testing.T, expect, actual T) {
	if expect != actual {
		t.Errorf("not equal, expect: %v, actual: %v", expect, actual)
	}
}

func asserEqualWithFunc[T any](t *testing.T, expect, actual T, equalFunc func(T, T) bool) {
	if !equalFunc(expect, actual) {
		t.Errorf("not equal, expect: %v, actual: %v", expect, actual)
	}
}
