package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ynishi/go-data-ope/_example/echotask"
)

const (
	preStr  = "1234"
	postStr = "abcd"
)

func TestEchoTaskWithContext(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := echotask.EchoReq{postStr}
	res := echotask.EchoRes{preStr}

	errChan := make(chan error, 1)
	go func() {
		errChan <- echoTaskHandler(ctx, req, &res, 0*time.Second)
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("Failed do task: %v", ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	}

	if !reflect.DeepEqual(postStr, res.Str) {
		t.Fatalf("not matched,\n want: %v,\n have: %v\n", postStr, res.Str)
	}

}

func TestEchoTaskWithCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := echotask.EchoReq{postStr}
	res := echotask.EchoRes{preStr}

	errChan := make(chan error, 1)
	go func() {
		errChan <- echoTaskHandler(ctx, req, &res, 10*time.Second)
	}()

	cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() == nil {
			t.Fatalf("Failed to catch ctx.Err()")
		}
	case err := <-errChan:
		t.Fatalf("Failed to cancel: %v\n", err)
	}

	if !reflect.DeepEqual(preStr, res.Str) {
		t.Fatalf("not matched,\n want: %v,\n have: %v\n", preStr, res.Str)
	}

}
