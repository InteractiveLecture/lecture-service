package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type Module struct {
	ID              string
	Description     string `valid:"utfletternumeric,required"`
	TopicID         string `valid:"utfletternumeric,required"`
	Recommendations []string
	Parents         []string
	Depth           uint
}

type Topic struct {
	ID          string
	Name        string `valid:"utfletternumeric,required"`
	Description string `valid:"utfletternumeric,required"`
	Officers    []string
	Assistants  []string
}

func NewModule(id, description, topicId string, depth uint, parents ...string) (*Module, error) {
	m := &Module{
		ID:          id,
		Description: description,
		TopicID:     topicId,
		Parents:     parents,
		Depth:       depth,
	}
	err := Validate(m)
	if err != nil {
		return nil, err
	}
	return m, nil
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
func NewTopic(id, name, description string, officers []string, assistants ...string) (*Topic, error) {

	t := &Topic{
		ID:          id,
		Name:        name,
		Description: description,
		Officers:    cleanDuplicates(officers),
		Assistants:  cleanDuplicates(assistants),
	}
	err := Validate(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func validationError() error {
	return errors.New("Violated validation")
}

func Validate(t interface{}) error {
	result, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}
	if !result {
		return validationError()
	}
	return nil
}
