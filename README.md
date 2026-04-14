# Charlotte Workflow Runner (PoC)

A lightweight proof-of-concept workflow execution engine written in Go. This project demonstrates a minimal approach to
orchestrating bash script workflows defined in YAML configuration files, with the unique ability to execute steps across
different environments seamlessly.

## Overview

This is a hobby/experimental project that provides:

* YAML-defined workflows: Describe your automation steps in a simple YAML format.
* Multi-environment Execution: Run steps locally, inside a Docker container, or on a Kubernetes cluster.
* Basic logic flow: Simple conditional logic and step sequencing.
* Input/output handling: Pass data between workflow steps.

Think of it as a stripped-down, self-hosted alternative to GitHub Actions — built for learning and experimentation 
rather than production use.

One of the core experiments in this project is the ability to define where a step runs without changing the script logic
itself.

* Local: Executes directly on the host machine.
* Docker: Spins up a temporary container for the step.
* Kubernetes: Deploys a Job to a configured K8s cluster.

## Example job YAML

```yaml
name: Test
description: Workflow with bash script steps
inputs:
  input_1:
    required: true
  input_2:
    required: true
    regexp: ^[A-Za-z0-9]+$
outputs:
  output_1:
    value: '{{ .StepOutputs.step_1.output_1 }}'
  output_2:
    value: '{{ .StepOutputs.step_1.output_2 }}'
steps:
  - type: shell
    name: Step 1
    id: step_1
    description: Simple test step
    script: |
      echo "Step1 Standard Output Message: {{ .Inputs.input_1 }}";
      >&2 echo "Step1 Standard Error Message: {{ .Inputs.input_2 }} ";
      echo -n "output1" > $OUTPUTS_DIR/output_1
      echo -n "output2" > $OUTPUTS_DIR/output_2
````

## Running test suite

    make test

## Building binary

    cd cmd/job
    go build .

## Running

    cd cmd/job
    ./job run-local -j ../../sample-files/job.yaml \
      -r /tmp/job-result.txt \
      -i ../../sample-files/job-inputs.json \
      --quiet
    cat /tmp/job-result.txt

Also, there are test files in the `pkg/job/runtime/local/tests` directory that can be used.

## Features

- [x] pipe stdout and stderr to files
- [x] environment (global and in-step)
- [x] variables
- [x] job inputs
- [x] step outputs
- [x] `continue_on_error`
- [x] values using golang templates
- [x] `if` - conditional steps (value templated, must equal to string `'true'`)
- [x] running step(s) on success
- [x] running step(s) on failure
- [x] running step(s) always
- [x] tmp directory for step outputs
- [x] gather job outputs 
- [x] write job outputs to json file
- [x] handle input: `--inputs`, `--job`, `--result` without aliases (and `--quiet`)
- [x] prepare sample yaml files - same as the test ones, so the test would just include them?

## Future features
- [ ] validation
- [ ] extract steps so that they can be included (include file with inputs) + proper validation for that
- [ ] pipelines
