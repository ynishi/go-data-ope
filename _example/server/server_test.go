package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	ect "github.com/ynishi/go-data-ope/_example/echotask"
)

var s *httptest.Server

const (
	origStr = `{"str":""}`
	testStr = `{"str":"abcd"}`

	wantValidate = `{"message":"validate succeed","echo_res":null}`
	wantPlan = `{"message":"plan succeed","echo_res":null}`
	wantDo = `{"message":"do succeed","echo_res":{"str":"abcd"}}`
	wantCheck = `{"message":"check succeed","echo_res":null}`
	wantBack = `{"message":"back succeed","echo_res":{"str":""}}`
)

func TestMain(m *testing.M) {

	//setup
	mux := http.NewServeMux()
	task = &ect.EchoTask{}

	mux.HandleFunc("/echo/validate", EchoValidateFunc)
	mux.HandleFunc("/echo/plan", EchoPlanFunc)
	mux.HandleFunc("/echo/do", EchoDoFunc)
	mux.HandleFunc("/echo/back", EchoBackFunc)
	mux.HandleFunc("/echo/check", EchoCheckFunc)

	s = httptest.NewServer(mux)

	retCode := m.Run()

	//teardown
	s.Close()

	os.Exit(retCode)
}

func HelperReq(path, str string) (resp *http.Response, err error) {
	if resp, err = http.Post(
		fmt.Sprintf("%s/echo/%s", s.URL, path),
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(str))); err != nil {
		return nil, err
	}
	return resp, nil
}

func TestValidate(t *testing.T) {

	var resp *http.Response
	var err error

	if resp, err = HelperReq("validate", testStr); err != nil {
		t.Fatal(err)
	}

	var buf []byte
	if buf,err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %v, body: %s", resp.Status, buf)
	}
	if !reflect.DeepEqual([]byte(wantValidate), buf) {
		t.Fatalf("not matched,\n want: %s,\n have: %s\n",wantValidate, buf)
	}
}


func TestPlan(t *testing.T) {
	var resp *http.Response
	var err error

	if resp, err = HelperReq("plan", testStr); err != nil {
		t.Fatal(err)
	}

	var buf []byte
	if buf,err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %v, body: %s", resp.Status, buf)
	}
	if !reflect.DeepEqual([]byte(wantPlan), buf) {
		t.Fatalf("not matched,\n want: %v,\n have: %v\n",wantPlan, buf)
	}
}

func TestDo(t *testing.T) {
	var resp *http.Response
	var err error

	if resp, err = HelperReq("do", testStr); err != nil {
		t.Fatal(err)
	}

	var buf []byte
	if buf,err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %v, body: %s", resp.Status, buf)
	}
	if !reflect.DeepEqual([]byte(wantDo), buf) {
		t.Fatalf("not matched,\n want: %v,\n have: %s\n",wantDo, buf)
	}
}

func TestCheck(t *testing.T) {
	var resp *http.Response
	var err error

	if resp, err = HelperReq("do", testStr); err != nil {
		t.Fatal(err)
	}

	if resp, err = HelperReq("check", testStr); err != nil {
		t.Fatal(err)
	}

	var buf []byte
	if buf,err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %v, body: %s", resp.Status, buf)
	}
	if !reflect.DeepEqual([]byte(wantCheck), buf) {
		t.Fatalf("not matched,\n want: %v,\n have: %s\n",wantCheck, buf)
	}
}

func TestBack(t *testing.T) {

	var resp *http.Response
	var err error

	if resp, err = HelperReq("do", origStr); err != nil {
		t.Fatal(err)
	}

	if resp, err = HelperReq("do", testStr); err != nil {
		t.Fatal(err)
	}

	if resp, err = HelperReq("back", origStr); err != nil {
		t.Fatal(err)
	}

	var buf []byte
	if buf,err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %v, body: %s", resp.Status, buf)
	}
	if !reflect.DeepEqual([]byte(wantBack), buf) {
		t.Fatalf("not matched,\n want: %v,\n have: %s\n",wantBack, buf)
	}
}
