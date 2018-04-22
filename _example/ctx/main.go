package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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

func main() {

	if len(os.Args) != 3 {
		log.Fatal("payload:$1, timeout:$2 is required")
	}

	d, err := time.ParseDuration(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	errChan := make(chan error, 1)
	req := ect.EchoReq{os.Args[1]}
	res := ect.EchoRes{}

	go func() {
		errChan <- echoTaskHandler(ctx, req, &res, d)
	}()

	select {
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			log.Fatal(err)
		}

	}

}
