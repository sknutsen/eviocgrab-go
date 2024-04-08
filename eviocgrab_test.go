package eviocgrabgo

import "testing"

func TestInit(t *testing.T) {
	_, err := Init()

	if err != nil {
		t.Error(err)
	}
}
