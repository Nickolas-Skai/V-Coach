package data

import (
	//	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
)

// Message represents a message sent by a user. messages can be either a string or a number at this point for testing purposes.

type Response struct {
	ID           int       `json:"id" db:"id"`
	SessionID    int       `json:"session_id" db:"session_id"`
	QuestionID   int       `json:"question_id" db:"question_id"`
	ResponseText *string   `json:"response_text,omitempty" db:"response_text"`
	AudioURL     *string   `json:"audio_url,omitempty" db:"audio_url"`
	Confidence   *int      `json:"confidence,omitempty" db:"confidence"`
	SubmittedAt  time.Time `json:"submitted_at" db:"submitted_at"`
}
type ResponseModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}

// answers to the questions
func (m *ResponseModel) ValidateResponse(response *Response) error {
	v := m.Validator
	if response.ResponseText != nil {
		v.Check(validator.NotBlank(*response.ResponseText), "response_text", "must be provided")
		v.Check(validator.MinLength(*response.ResponseText, 1), "response_text", "must be at least 1 characters long")
		v.Check(validator.MaxLength(*response.ResponseText, 1000), "response_text", "must not exceed 1000 characters")
	}
	return fmt.Errorf("validation errors: %v", v.Errors)
}

func (m *ResponseModel) Insert(response *Response) error {
	// Insert the response into the database
	query := `INSERT INTO responses (session_id, question_id, response_text, audio_url, confidence, submitted_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	args := []interface{}{response.SessionID, response.QuestionID, response.ResponseText, response.AudioURL, response.Confidence, time.Now()}
	err := m.DB.QueryRow(query, args...).Scan(&response.ID)
	if err != nil {
		return err
	}
	return nil
}
func (m *ResponseModel) GetByID(id int) (*Response, error) {
	// Retrieve the response from the database by ID
	query := `SELECT id, session_id, question_id, response_text, audio_url, confidence, submitted_at FROM responses WHERE id = $1`
	row := m.DB.QueryRow(query, id)

	var response Response
	err := row.Scan(&response.ID, &response.SessionID, &response.QuestionID, &response.ResponseText, &response.AudioURL, &response.Confidence, &response.SubmittedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No response found with the given ID
		}
		return nil, err // Some other error occurred
	}
	return &response, nil
}
func (m *ResponseModel) GetBySessionID(sessionID int) ([]Response, error) {
	// Retrieve all responses for a given session ID
	query := `SELECT id, session_id, question_id, response_text, audio_url, confidence, submitted_at FROM responses WHERE session_id = $1`
	rows, err := m.DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []Response
	for rows.Next() {
		var response Response
		err := rows.Scan(&response.ID, &response.SessionID, &response.QuestionID, &response.ResponseText, &response.AudioURL, &response.Confidence, &response.SubmittedAt)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}
