package job

import structvalidator "github.com/mikolajgasior/struct-validator"

func (j *Job) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(j, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
