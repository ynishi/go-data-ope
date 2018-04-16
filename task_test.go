package dataope

import "testing"

type testTask struct{}

type testRequset struct{}

type testPlanResult struct{}

type testDoResult struct{}

type testBackResult struct{}

type testMonitorResult struct{}

func (tt *testTask) Validate(v interface{}) error {
	return nil
}

func (tt *testTask) Plan(req interface{}, v interface{}) error {
	return nil
}

func (tt *testTask) Do(req interface{},  v interface{}) error {
	return nil
}
func (tt *testTask) Back(req interface{},  v interface{}) error {
	return nil
}
func (tt *testTask) Check(req interface{}) error {
	return nil
}
func (tt *testTask) Monitor(req interface{},  v interface{}) error {
	return nil
}

func TestTasker(t *testing.T) {
	task := &testTask{}
	if  _, ok := interface{}(task).(Tasker); !ok {
		t.Errorf("failed to implement Tasker")
	}
	request := testRequset{}
	if err := task.Validate(&request); err != nil {
		t.Errorf("failed to validate: %v", err)
	}
	planResult := testPlanResult{}
	if err := task.Plan(&request, &planResult); err != nil {
		t.Errorf("failed to call plan: %v", err)
	}
	doResult := testDoResult{}
	if err := task.Do(&request, &doResult); err != nil {
		t.Errorf("failed to call do: %v", err)
	}
	backResult := testBackResult{}
	if err := task.Back(&request, &backResult); err != nil {
		t.Errorf("failed to call back: %v", err)
	}
	if err := task.Check(&request); err != nil {
		t.Errorf("failed to call check: %v", err)
	}
	monitorResult := testMonitorResult{}
	if err := task.Monitor(&request, &monitorResult); err != nil {
		t.Errorf("failed to call monitor: %v", err)
	}

}
