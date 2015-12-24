package datamapper

import (
	"errors"
	"fmt"
)

func (r *DataMapper) CompleteExercise(exerciseId, userId string) error {
	_, err := r.db.Exec(`SELECT complete_exercise($1,$2)`, exerciseId, userId)
	return err
}

func (r *DataMapper) CompleteTask(taskId, userId string) error {
	_, err := r.db.Exec(`SELECT complete_task($1,$2)`, taskId, userId)
	return err
}

func (r *DataMapper) GetHint(hintId, userId string) ([]byte, error) {
	result, err := r.queryIntoBytes(`SELECT get_hint($1,$2)`, userId, hintId)
	switch {
	case err != nil:
		return nil, err
	case len(result) == 0:
		return nil, PaymentRequiredError{fmt.Sprintf("User %s needs to pay for hint %s", userId, hintId)}
	}
	return result, nil
}

func (r *DataMapper) PurchaseHint(hintId, userId string) error {
	row := r.db.QueryRow("SELECT purchase_hint($1,$2)", hintId, userId)
	var purchaseResult int
	err := row.Scan(&purchaseResult)
	if err != nil {
		return err
	}
	switch {
	case purchaseResult == 0:
		return nil
	case purchaseResult == 1:
		return InsufficientPointsError{"User balance is not sufficient for this hint."}
	case purchaseResult == 2:
		return AlreadyPurchasedError{fmt.Sprintf("User %s has already purchased hint %s", userId, hintId)}
	case purchaseResult == 3:
		return HintNotFoundError{fmt.Sprintf("Could not find hint with id %s", hintId)}
	default:
		return errors.New("Unknown transaction status")
	}
}

type AlreadyPurchasedError struct {
	Message string
}

type InsufficientPointsError struct {
	Message string
}

type HintNotFoundError struct {
	Message string
}

type PaymentRequiredError struct {
	Message string
}

func (e AlreadyPurchasedError) Error() string {
	return e.Message
}

func (e PaymentRequiredError) Error() string {
	return e.Message
}

func (e InsufficientPointsError) Error() string {
	return e.Message
}

func (e HintNotFoundError) Error() string {
	return e.Message
}
