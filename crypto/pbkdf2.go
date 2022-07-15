package crypto

import (
	"golang.org/x/crypto/pbkdf2"
	"hash"
)

type PBKDF2Config struct {
	Password string `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Salt     string `yaml:"salt" mapstructure:"salt" json:"salt,omitempty" gorm:"column:salt" bson:"salt,omitempty" dynamodbav:"salt,omitempty" firestore:"salt,omitempty"`
	Iter     int    `yaml:"iter" mapstructure:"iter" json:"iter,omitempty" gorm:"column:iter" bson:"iter,omitempty" dynamodbav:"iter,omitempty" firestore:"iter,omitempty"`
	KeyLen   int    `yaml:"key_len" mapstructure:"key_len" json:"keyLen,omitempty" gorm:"column:keylen" bson:"keyLen,omitempty" dynamodbav:"keyLen,omitempty" firestore:"keyLen,omitempty"`
}

func NewPBKDF2Key(c PBKDF2Config, h func() hash.Hash) []byte {
	return pbkdf2.Key([]byte(c.Password), []byte(c.Salt), c.Iter, c.KeyLen, h)
}
