// SPDX-License-Identifier: AGPL-3.0-only

package job

import (
	"os/exec"

	"github.com/brainupdaters/drlm-agent/models"

	drlm "github.com/brainupdaters/drlm-common/pkg/proto"
	log "github.com/sirupsen/logrus"
)

var jobs = map[uint32]*models.Job{}

// Run starts a job
func Run(c chan models.JobUpdate, id uint32, name, config, target string) {
	j := &models.Job{
		ID:     id,
		Name:   name,
		Config: config,
		Target: target,
		Cmd:    exec.Command(name, "-config", config, "-target", target),
	}

	jobs[j.ID] = j

	if out, err := j.Cmd.CombinedOutput(); err != nil {
		delete(jobs, j.ID)

		c <- models.JobUpdate{
			JobID:  j.ID,
			Status: drlm.JobStatus_JOB_STATUS_FAILED,
			Info:   string(out),
		}
		log.Errorf("error running '%s': %v: %s", j.Name, err, out)

		return
	}

	delete(jobs, j.ID)
	c <- models.JobUpdate{
		JobID:  j.ID,
		Status: drlm.JobStatus_JOB_STATUS_FINISHED,
	}
}

// Cancel kills a running job
func Cancel(id uint32) {
	if j, ok := jobs[id]; ok {
		if err := j.Cmd.Process.Kill(); err != nil {
			log.Errorf("error cancelling '%s': %v", j.Name, err)
		}
	}
}
