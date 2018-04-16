# go-data-ope
simple data operation with go.
* Simple data operation/flow framework.
* It can make combinations of each data operation tasks.
* Write tasks with Go(language integrate style).
* It have rest api interface to manage and monitor task execution.
* kick a task via rest api, task execute concurrent.

## interface methods of a task
* validete -- check input or precondition.
* plan -- plan execution, no effect.
* do -- main processing of this task.
* cancel -- cancel main processing.
* back -- fail back main processing.
* check -- check output or effect.
* monitor -- get task execution status.
