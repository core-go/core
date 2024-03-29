package diff

type StatusDiffConfig struct {
	Status *StatusConfig `yaml:"status" mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Config DiffConfig    `yaml:"config" mapstructure:"config" json:"config" gorm:"column:config" bson:"config" dynamodbav:"config" firestore:"config"`
}

type StatusConfig struct {
	NotFound     int `yaml:"not_found" mapstructure:"not_found" json:"notFound" gorm:"column:notfound" bson:"notFound" dynamodbav:"notFound" firestore:"notFound"`
	Success      int `yaml:"success" mapstructure:"success" json:"success" gorm:"column:success" bson:"success" dynamodbav:"success" firestore:"success"`
	VersionError int `yaml:"version_error" mapstructure:"version_error" json:"versionError" gorm:"column:versionerror" bson:"versionError" dynamodbav:"versionError" firestore:"versionError"`
	Error        int `yaml:"error" mapstructure:"error" json:"error" gorm:"column:error" bson:"error" dynamodbav:"error" firestore:"error"`
}

func InitializeStatus(status *StatusConfig) StatusConfig {
	var s StatusConfig
	if status != nil {
		s.NotFound = status.NotFound
		s.Success = status.Success
		s.VersionError = status.VersionError
		s.Error = status.Error
	}
	if s.NotFound == 0 && s.Success == 0 && s.VersionError == 0 && s.Error == 0 {
		s.Success = 1
	}
	if s.NotFound == 0 && s.Success == 1 && s.VersionError == 0 && s.Error == 0 {
		s.VersionError = 2
		s.Error = 4
	}
	return s
}
