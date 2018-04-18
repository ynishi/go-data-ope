package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	gdo "github.com/ynishi/go-data-ope"
	ect "github.com/ynishi/go-data-ope/_example/echotask"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("payload:$1 is required")
	}
	echoReq := ect.EchoReq{os.Args[1]}

	var task gdo.Tasker
	task = &ect.EchoTask{}

	if err := task.Validate(echoReq); err != nil {
		log.Fatalf("failed to validate: %v\n", err)
	}

	echoRes := ect.EchoRes{""}
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
