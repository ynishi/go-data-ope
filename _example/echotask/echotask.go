package echotask

import (
	"errors"
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
