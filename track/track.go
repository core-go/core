package track

import "time"

type Tracker struct {
	Id  string
	User  string
	Time  string
	Status  string
	Version string
}
func UseMap(opts...string) func(string, string, *string)map[string]interface{} {
	t := NewTracker(opts...)
	return t.Create
}
func NewTracker(opts...string) *Tracker {
	var id, user, time, status, version string
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
	if len(opts) > 4 {
		version = opts[4]
	} else {
		version = "version"
	}
	return &Tracker{Id: id, User: user, Time: time, Status: status, Version: version}
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
func (t *Tracker) CreateMap(id string, status string, version int32, userId *string) map[string]interface{} {
	obj := make(map[string]interface{})
	obj[t.Id] = id
	obj[t.Status] = status
	obj[t.Version] = version
	if userId != nil {
		obj[t.Time] = time.Now()
		obj[t.User] = *userId
	}
	return obj
}
