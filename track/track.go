package track

import "time"

type Tracker struct {
	Id  string
	User  string
	Time  string
	Status  string
}
func UseMap(opts...string) func(string, string, *string)map[string]interface{} {
	t := NewTracker(opts...)
	return t.Create
}
func NewTracker(opts...string) *Tracker {
	var id, user, time, status string
	if len(opts) > 0 {
		id = opts[0]
	} else {
		id = "id"
	}
	if len(opts) > 1 {
		user = opts[1]
	} else {
		user = "updatedBy"
	}
	if len(opts) > 2 {
		time = opts[2]
	} else {
		time = "updatedAt"
	}
	if len(opts) > 3 {
		status = opts[3]
	} else {
		status = "status"
	}
	return &Tracker{Id: id, User: user, Time: time, Status: status}
}
func (t *Tracker) Create(id string, status string, userId *string) map[string]interface{} {
	obj := make(map[string]interface{})
	obj[t.Id] = id
	obj[t.Status] = status
	if userId != nil {
		obj[t.Time] = time.Now()
		obj[t.User] = *userId
	}
	return obj
}
