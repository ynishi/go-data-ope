package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	gdo "github.com/ynishi/go-data-ope"
	ect "github.com/ynishi/go-data-ope/_example/echotask"
)

func echoTaskHandler(ctx context.Context, echoReq ect.EchoReq, echoRes *ect.EchoRes, dur time.Duration) error {

	var task gdo.Tasker
	task = &ect.EchoTask{}

	if err := task.Validate(echoReq); err != nil {
		log.Fatalf("failed to validate: %v\n", err)
	}

	if err := task.Plan(echoReq, echoRes); err != nil {
		log.Fatalf("failed to plan: %v\n", err)
	}
	buf, err := json.Marshal(echoRes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("{\"plan_result\": %v}\n", string(buf))

	errChan := make(chan error, 1)

	go func() {
		errChan <- task.Do(echoReq, echoRes)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("canceled")
		if ctx.Err() != nil {
			if err = task.Back(echoReq, echoRes); err != nil {
				return errors.New(fmt.Sprintf("failed to back: %v\n", err))
			}
			log.Printf("succeed to back: %v\n", err)
		}
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			if err = task.Back(echoReq, echoRes); err != nil {
				return errors.New(fmt.Sprintf("failed to back: %v\n", err))
			}
			log.Printf("succeed to back: %v\n", err)
			return err
		}

		buf, err = json.Marshal(echoRes)
		if err != nil {
			return err
		}
		fmt.Printf("{\"do_result\": %v}\n", string(buf))

		if err = task.Check(echoReq); err != nil {
			log.Printf("failed to do, try to back: %v\n", err)
			if err = task.Back(echoReq, echoRes); err != nil {
				return errors.New(fmt.Sprintf("failed to back: %v\n", err))
			}
			log.Printf("succeed to back: %v\n", err)
			return err
		}
		fmt.Println(`{"check_result": "succeed"}`)
		return nil
	}
	return nil
}

var task = &ect.EchoTask{}

func main() {

	http.HandleFunc("/echo/validate", EchoValidateFunc)
	http.HandleFunc("/echo/plan", EchoPlanFunc)
	http.HandleFunc("/echo/do", EchoDoFunc)
	http.HandleFunc("/echo/back", EchoBackFunc)
	http.HandleFunc("/echo/check", EchoCheckFunc)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type EchoResponse struct {
	Message string       `json:"message"`
	Echo    *ect.EchoRes `json:"echo_res"`
}

func EchoValidateFunc(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed read body"))
		return
	}

	echoReq := ect.EchoReq{}
	err = json.Unmarshal(buf, &echoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed unmarshal echo req"))
		return
	}

	if err := task.Validate(echoReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validate error"))
		return
	}

	echoResponse := EchoResponse{
		Message: "validate succeed",
		Echo:    nil,
	}

	buf, err = json.Marshal(echoResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return
}

func EchoPlanFunc(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoReq := ect.EchoReq{}
	err = json.Unmarshal(buf, &echoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoRes := ect.EchoRes{}
	if err := task.Plan(echoReq, &echoRes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("plan error"))
		return
	}

	echoResponse := EchoResponse{
		Message: "plan succeed",
		Echo:    nil,
	}

	buf, err = json.Marshal(echoResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return
}

func EchoDoFunc(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoReq := ect.EchoReq{}
	err = json.Unmarshal(buf, &echoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoRes := &ect.EchoRes{}
	if err := task.Do(echoReq, echoRes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("do error"))
		return
	}

	echoRes.Str = echoReq.Str
	echoResponse := EchoResponse{
		Message: "do succeed",
		Echo:    echoRes,
	}

	buf, err = json.Marshal(echoResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return
}

func EchoBackFunc(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoReq := ect.EchoReq{}
	err = json.Unmarshal(buf, &echoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoRes := ect.EchoRes{}
	if err := task.Back(echoReq, &echoRes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("back error"))
		return
	}

	echoRes.Str = echoReq.Str
	echoResponse := EchoResponse{
		Message: "back succeed",
		Echo:    &echoRes,
	}

	buf, err = json.Marshal(echoResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return
}

func EchoCheckFunc(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	echoReq := ect.EchoReq{}
	err = json.Unmarshal(buf, &echoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	if err := task.Check(echoReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("check error"))
		return
	}

	echoResponse := EchoResponse{
		Message: "check succeed",
		Echo:    nil,
	}

	buf, err = json.Marshal(echoResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return
}
