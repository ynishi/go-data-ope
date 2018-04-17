package main

import (
	"os/exec"
	"reflect"
	"testing"
)

const (
	expected = `{"plan_result": {"Str":"abcd"}}
{"do_result": {"Str":"abcd"}}
{"check_result": "succeed"}`
)

func TestMainExec(t *testing.T) {

	out, err := exec.Command("go","run","echoope.go", "abcd").Output()
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(expected,out) {
		t.Errorf("not matched:\n want: %v,\n have: %s\n", expected, out)
	}

}
