package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type Topic struct {
	name        string `valid:"utfletternumeric,required"`
	description string `valid:"utfletternumeric,required"`
	officers    map[string]bool
	assistants  map[string]bool
}

func toSet(slice []string) map[string]bool {
	var result = make(map[string]bool)
	for _, v := range slice {
		result[v] = true
	}
	return result
}
func NewTopic(name, description string, officers []string, assistants ...string) (*Topic, error) {
	t := &Topic{
		name:        name,
		description: description,
		officers:    toSet(officers),
		assistants:  toSet(assistants),
	}
	err := t.Validate()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Topic) AddOfficers(officers ...string) {
	for _, o := range officers {
		t.officers[o] = true
	}
}

func (t *Topic) RemoveOfficers(officers ...string) error {
	tmp := t.officers
	for _, o := range officers {
		delete(tmp, o)
	}
	if len(tmp) > 0 {
		t.officers = tmp

		return nil
	}
	return errors.New("Must have at least one officer")
}

func (t *Topic) Name() string {
	return t.name
}

func (t *Topic) SetName(name string) error {
	if !govalidator.IsUTFLetterNumeric(name) {
		return validationError()
	}
	t.name = name
	return nil
}

func (t *Topic) SetDescription(desc string) error {
	if !govalidator.IsUTFLetterNumeric(desc) {
		return validationError()
	}
	t.description = desc
	return nil
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
