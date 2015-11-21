package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/asaskevich/govalidator"
)

type Module struct {
	ID              string `bson:"_id,omitempty"`
	TopicID         string
	Recommendations []string
	Path            string
}

type Topic struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `valid:"utfletternumeric,required"`
	Description string `valid:"utfletternumeric,required"`
	Officers    []string
	Assistants  []string
}

func cleanDuplicates(slice []string) []string {
	var set = make(map[string]bool)
	for _, v := range slice {
		set[v] = true
	}
	var result = make([]string, len(set))
	i := 0
	for k, _ := range set {
		result[i] = k
		i = i + 1
	}
	return result
}
func NewTopic(name, description string, officers []string, assistants ...string) (*Topic, error) {

	t := &Topic{
		ID:          bson.NewObjectId().Hex(),
		Name:        name,
		Description: description,
		Officers:    cleanDuplicates(officers),
		Assistants:  cleanDuplicates(assistants),
	}
	err := t.Validate()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func validationError() error {
	return errors.New("Violated validation")
}

func (t *Topic) Validate() error {
	result, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}
	if !result {
		return validationError()
	}
	return nil
}
