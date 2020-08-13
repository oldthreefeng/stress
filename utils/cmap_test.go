// "github.com/orcaman/concurrent-map"
package utils

import (
	"testing"
)

type Animal struct {
	name string
}

func TestNew(t *testing.T) {
	m := New()

	elephant := Animal{"elephant"}
	m.Set("a", elephant)
	if val, ok := m.Get("a") ; !ok {
		t.Error(val)
	}

}