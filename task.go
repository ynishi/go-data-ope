package dataope

import "context"

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

type CtxTasker interface {
	Preparer
	Runner
	Rollbacker
	Commiter
	Stater
}

type Preparer interface {
	Prepare(context.Context, interface{}, interface{}) error
}

type Runner interface {
	Run(context.Context, interface{}, interface{}) error
}

type Rollbacker interface {
	Rollback(context.Context, interface{}, interface{}) error
}

type Commiter interface {
	Commit(context.Context, interface{}, interface{}) error
}

type Stater interface {
	Stat(context.Context, interface{}, interface{}) error
}

type DefaultServer struct {
	CtxTask CtxTasker
}

func NewDefaultServer(ctxTask CtxTasker) (dc *DefaultServer, err error) {

	return &DefaultServer{ctxTask}, nil
}

func (dc *DefaultServer) Run(ctx context.Context, req interface{}, res interface{}) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		if err = dc.CtxTask.Prepare(ctx, req, res); err != nil {
			errChan <- err
		}
		errChan <- dc.CtxTask.Run(ctx, req, res)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		commitErrChan := make(chan error, 1)
		go func() {
			commitErrChan <- dc.CtxTask.Commit(ctx, req, res)
		}()

		select {
		case err := <-commitErrChan:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return err
		}
		if err != nil {
			rollbackErrChan := make(chan error, 1)
			go func() {
				rollbackErrChan <- dc.CtxTask.Rollback(ctx, req, res)
			}()

			select {
			case err := <-rollbackErrChan:
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return err
			}
			return err
		}
	}
	return nil
}
