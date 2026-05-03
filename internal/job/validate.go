package job

import structvalidator "github.com/keenbytes/struct-validator"

func (j *Job) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(j, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
