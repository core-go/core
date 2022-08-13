package sql

import "time"

type Statement struct {
	Query  string        `yaml:"query" mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	Params []interface{} `yaml:"params" mapstructure:"params" json:"params,omitempty" gorm:"column:params" bson:"params,omitempty" dynamodbav:"params,omitempty" firestore:"params,omitempty"`
}

type JStatement struct {
	Query  string        `yaml:"query" mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	Params []interface{} `yaml:"params" mapstructure:"params" json:"params,omitempty" gorm:"column:params" bson:"params,omitempty" dynamodbav:"params,omitempty" firestore:"params,omitempty"`
	Dates  []int         `yaml:"dates" mapstructure:"dates" json:"dates,omitempty" gorm:"column:dates" bson:"dates,omitempty" dynamodbav:"dates,omitempty" firestore:"dates,omitempty"`
}
func BuildStatement(query string, values ...interface{}) *JStatement {
	stm := JStatement{Query: query}
	l := len(values)
	if l > 0 {
		ag2 := make([]interface{}, 0)
		dates := make([]int, 0)
		for i := 0; i < l; i++ {
			arg := values[i]
			if _, ok := arg.(time.Time); ok {
				dates = append(dates, i)
			} else if _, ok := arg.(*time.Time); ok {
				dates = append(dates, i)
			}
			ag2 = append(ag2, values[i])
		}
		stm.Params = ag2
		if len(dates) > 0 {
			stm.Dates = dates
		}
	}
	return &stm
}
func BuildJStatements(sts ...Statement) []JStatement {
	b := make([]JStatement, 0)
	if sts == nil || len(sts) == 0 {
		return b
	}
	for _, s := range sts {
		j := JStatement{Query: s.Query}
		if s.Params != nil && len(s.Params) > 0 {
			j.Params = s.Params
			j.Dates = ToDates(s.Params)
		}
		b = append(b, j)
	}
	return b
}
