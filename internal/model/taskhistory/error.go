package taskhistory

import "errors"

var ErrNoChanges = errors.New("task has no changes")
var ErrTaskNotFound = errors.New("task not found")
