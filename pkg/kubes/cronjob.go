package kubes

import (
	"regexp"
	"strings"
)

const (
	scheduleDayRegex  = `^\w{3}$`
	scheduleDateRegex = `^(\d+|\*)-(\d+|\*)-(\d+|\*)$`
	scheduleTimeRegex = `^(\d+|\*):(\d+|\*):?(\d+|\*)?$`
)

var (
	daysOfWeek = map[string]string{
		"Sun": "0", "Mon": "1", "Tue": "2", "Wed": "3",
		"Thu": "4", "Fri": "5", "Sat": "6",
	}
)

type CronJob struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name        string      `json:"name"`
		Annotations Annotations `json:"annotations" yaml:"annotations,omitempty"`
	} `json:"metadata"`
	Spec struct {
		ConcurrencyPolicy string `json:"concurrencyPolicy" yaml:"concurrencyPolicy"`
		Schedule          string `json:"string"`
		JobTemplate       struct {
			Spec struct {
				Template struct {
					Spec struct {
						RestartPolicy string      `json:"restartPolicy" yaml:"restartPolicy"`
						Containers    []Container `json:"containers"`
					}
				} `json:"template"`
			} `json:"spec"`
		} `json:"jobTemplate" yaml:"jobTemplate"`
	} `json:"spec"`
}

type Annotations map[string]string

func stripZeroes(n string) string {
	for len(n) > 1 && n[0] == '0' {
		n = n[1:]
	}

	return n
}

func parseSchedule(schedule string) string {
	scheduleParts := strings.Split(schedule, " ")
	dayOfWeek := "*"
	dayOfMonth := "*"
	month := "*"
	hour := "*"
	minute := "*"

	rDay := regexp.MustCompile(scheduleDayRegex)
	rDate := regexp.MustCompile(scheduleDateRegex)
	rTime := regexp.MustCompile(scheduleTimeRegex)
	for _, part := range scheduleParts {
		if rDay.MatchString(part) {
			if day, ok := daysOfWeek[part]; ok {
				dayOfWeek = day
			}
		} else if rDate.MatchString(part) {
			dateParts := strings.Split(part, "-")
			month = stripZeroes(dateParts[1])
			dayOfMonth = stripZeroes(dateParts[2])
		} else if rTime.MatchString(part) {
			timeParts := strings.Split(part, ":")
			hour = stripZeroes(timeParts[0])
			minute = stripZeroes(timeParts[1])
		}
	}

	return strings.Join(append([]string{}, minute, hour, dayOfMonth, month, dayOfWeek), " ")
}

func NewCronJob(name, schedule, image string, command []string, env map[string]string, annotations Annotations) *CronJob {
	cronJob := &CronJob{
		APIVersion: "batch/v1beta1",
		Kind:       "CronJob",
	}
	cronJob.Metadata.Name = name
	cronJob.Metadata.Annotations = annotations
	cronJob.Spec.ConcurrencyPolicy = "Forbid"
	cronJob.Spec.Schedule = parseSchedule(schedule)
	cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers = append(
		cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers,
		Container{
			Name:    name,
			Image:   image,
			Command: command,
			Env:     newEnv(env),
		},
	)
	cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = "OnFailure"

	return cronJob
}
