package search

type SearchResultConfig struct {
	Results       string `yaml:"results" mapstructure:"results" json:"results,omitempty" gorm:"column:results" bson:"results,omitempty" dynamodbav:"results,omitempty" firestore:"results,omitempty"`
	Total         string `yaml:"total" mapstructure:"total" json:"total,omitempty" gorm:"column:total" bson:"total,omitempty" dynamodbav:"total,omitempty" firestore:"total,omitempty"`
	Next          string `yaml:"next" mapstructure:"next" json:"next,omitempty" gorm:"column:next" bson:"next,omitempty" dynamodbav:"next,omitempty" firestore:"next,omitempty"`
	// LastPage      string `yaml:"last_page" mapstructure:"last_page" json:"lastPage,omitempty" gorm:"column:lastpage" bson:"lastPage,omitempty" dynamodbav:"lastPage,omitempty" firestore:"lastPage,omitempty"`
	PageIndex     string `yaml:"page_index" mapstructure:"page_index" json:"pageIndex,omitempty" gorm:"column:pageindex" bson:"pageIndex,omitempty" dynamodbav:"pageIndex,omitempty" firestore:"pageIndex,omitempty"`
	PageSize      string `yaml:"page_size" mapstructure:"page_size" json:"pageSize,omitempty" gorm:"column:pagesize" bson:"pageSize,omitempty" dynamodbav:"pageSize,omitempty" firestore:"pageSize,omitempty"`
	FirstPageSize string `yaml:"first_page_size" mapstructure:"first_page_size" json:"firstPageSize,omitempty" gorm:"column:firstpagesize" bson:"firstPageSize,omitempty" dynamodbav:"firstPageSize,omitempty" firestore:"firstPageSize,omitempty"`
	NextPageToken string `yaml:"next_page_token mapstructure:"next_page_token" json:"nextPageToken,omitempty" gorm:"column:nextPageToken" bson:"nextPageToken,omitempty" dynamodbav:"nextPageToken,omitempty" firestore:"nextPageToken,omitempty"`
}
