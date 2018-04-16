package dataope

type Tasker interface {
	Validator
	Planner
	Executor
	Backer
	Checker
	Monitoror
}

type Validator interface {
	Validate(interface{}) error
}

type Planner interface {
	Plan(interface{}, interface{}) error
}

type Executor interface {
	Do(interface{}, interface{}) error
}

type Backer interface {
	Back(interface{}, interface{}) error
}

type Checker interface {
	Check(interface{}) error
}

type Monitoror interface {
	Monitor(interface{}, interface{}) error
}