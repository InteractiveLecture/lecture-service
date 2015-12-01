package datamapper

import "fmt"

func (r *DataMapper) CompleteExercise(exerciseId string, userId string) error {
	_, err := r.db.Exec(`SELECT complete_exercise($1)`, id)

	return err
}

func (r *DataMapper) GetHint(hintId, userId string) ([]byte, error) {
	result, err := r.db.QueryRow(`SELECT check_needs_hint_payment($1,$2)`, id)
	if err != nil {
		return nil, err
	}
	var needsPayment bool
	err = result.Scan(needsPayment)
	if err != nil {
		return nil, err
	}
	if needsPayment {
		return nil, PaymentRequiredError{fmt.Sprintf("User %s needs to pay for hint %s", userId, hintId)}
	}
	return rowToBytes(r.db.QueryRow(`SELECT get_hint($1)`, hintId))
}

type PaymentRequiredError struct {
	Message string
}

func (e *PaymentRequiredError) Error() string {
	return e.Message
}
