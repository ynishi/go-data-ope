package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	gdo "github.com/ynishi/go-data-ope"
)

type EchoTask struct {
	Original string
	Output   string
}

type EchoReq struct {
	Str string `json:"str"`
}

type EchoRes struct {
	Str string `json:"str"`
}

func (et *EchoTask) Validate(req interface{}) error {
	if _, ok := req.(EchoReq); !ok {
		return errors.New("request type is not EchoReq.")
	}
	return nil
}

func (et *EchoTask) Plan(req interface{}, v interface{}) error {
	echoReq, err := v2EchoReq(req)
	if err != nil {
		return err
	}
	echoRes, err := v2EchoRes(v)
	if err != nil {
		return err
	}
	echoRes.Str = echoReq.Str
	return nil
}

func (et *EchoTask) Do(req interface{}, v interface{}) error {
	echoReq, err := v2EchoReq(req)
	if err != nil {
		return err
	}
	et.Original = et.Output
	et.Output = echoReq.Str
	v = echoReq.Str
	return nil
}
func (et *EchoTask) Back(req interface{}, v interface{}) error {
	echoReq, err := v2EchoReq(req)
	if err != nil {
		return err
	}
	if echoReq.Str != et.Original {
		return errors.New("invalid rollback data: original " + et.Original)
	}
	et.Original = ""
	et.Output = et.Original
	v = et.Output
	return nil
}
func (et *EchoTask) Check(req interface{}) error {
	echoReq, err := v2EchoReq(req)
	if err != nil {
		return err
	}
	if echoReq.Str != et.Output {
		return errors.New("failed to check: " + et.Output)
	}
	return nil
}
func (et *EchoTask) Monitor(req interface{}, v interface{}) error {
	return nil
}

func v2EchoReq(v interface{}) (*EchoReq, error) {
	var e EchoReq
	var ok bool
	if e, ok = v.(EchoReq); !ok {
		return nil, errors.New("request type is not EchoReq.")
	}
	return &e, nil
}

func v2EchoRes(v interface{}) (*EchoRes, error) {
	var e *EchoRes
	var ok bool
	if e, ok = v.(*EchoRes); !ok {
		return nil, errors.New("request type is not EchoRes.")
	}
	return e, nil
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("payload:$1 is required")
	}
	echoReq := EchoReq{os.Args[1]}

	var task gdo.Tasker
	task = &EchoTask{}

	if err := task.Validate(echoReq); err != nil {
		log.Fatalf("failed to validate: %v\n", err)
	}

	echoRes := EchoRes{""}
	if err := task.Plan(echoReq, &echoRes); err != nil {
		log.Fatalf("failed to plan: %v\n", err)
	}
	buf, err := json.Marshal(echoRes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("{\"plan_result\": %v}\n", string(buf))

	if err = task.Do(echoReq, &echoRes); err != nil {
		log.Fatalf("failed to do: %v\n", err)
	}
	buf, err = json.Marshal(echoRes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("{\"do_result\": %v}\n", string(buf))

	if err = task.Check(echoReq); err != nil {
		log.Printf("failed to do, try to back: %v\n", err)
		if err = task.Back(echoReq, echoRes); err != nil {
			log.Fatalf("failed to back: %v\n", err)
		}
		log.Printf("succeed to back: %v\n", err)
		return
	}
	fmt.Println(`{"check_result": "succeed"}`)

}
