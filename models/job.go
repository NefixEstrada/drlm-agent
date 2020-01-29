// SPDX-License-Identifier: AGPL-3.0-only

package models

import "os/exec"

import drlm "github.com/brainupdaters/drlm-common/pkg/proto"

// Job is an individual job
type Job struct {
	ID     uint32
	Name   string
	Config string
	Target string

	Cmd *exec.Cmd
}

// JobUpdate is the update that is going to be sent to the DRLM Core when there's an update in a job (e.g. the job has failed, finished, new output...)
type JobUpdate struct {
	JobID  uint32
	Status drlm.JobStatus
	Info   string
}
