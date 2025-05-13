package data

import (
	//	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
)

// Question represents a question in the system.
type QuestionData struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Type     string   `json:"type"`
	AudioURL string   `json:"audio_url"`
	ImageURL string   `json:"image_url"`
	Options  []string `json:"options"`
	Required bool     `json:"required"`
}
type QuestionModel struct {
	DB        *sql.DB
	Validator *validator.Validator
	ID        int
	Text      string
	Type      string
	Options   []string
	Required  bool
}
type InterviewResponse struct {
	ID                    int
	QuestionID            int
	SessionID             int
	ConfidenceRating      int
	SubmittedAt           string
	Text                  string
	Questiontype          string
	Options               []string
	AudioURL              string
	Confidence            int
	Answer                string
	QuestionData          *QuestionData
	QuestionDataErrors    map[string]string
	QuestionDataFormData  map[string]string
	AllowConfidenceRating bool
	Required              bool
	QuestionDataID        int
	QuestionDataText      string
	QuestionDataType      string
	QuestionDataOptions   []string
	QuestionDataRequired  bool
	QuestionDataModel     *QuestionModel
	QuestionDataDB        *QuestionModel
	TotalQuestions        int
	CurrentIndex          int
}

type InterviewResponseModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}

func NewInterviewResponseModel(db *sql.DB) *InterviewResponseModel {
	return &InterviewResponseModel{
		DB:        db,
		Validator: validator.NewValidator()}

}

func NewQuestionModel(db *sql.DB) *QuestionModel {
	return &QuestionModel{
		DB:        db,
		Validator: validator.NewValidator(),
	}
}

// function to validate the question response

func (m *QuestionModel) ValidateQuestionData(questionData *QuestionData) error {
	// Validate the question text
	if questionData.Text == "" {
		m.Validator.AddError("text", "Question text cannot be empty")
		return fmt.Errorf("validation failed: question text cannot be empty")
	}

	// Validate the question type
	if !validator.IsValidQuestionType(questionData.Type) {
		m.Validator.AddError("type", "Invalid question type")
		return fmt.Errorf("validation failed: invalid question type")
	}

	// Validate the options if the question type requires it
	if questionData.Type == "checkbox" || questionData.Type == "radio" || questionData.Type == "scale" {
		if len(questionData.Options) == 0 {
			m.Validator.AddError("options", "Options cannot be empty for this question type")
			return fmt.Errorf("validation failed: options cannot be empty for this question type")
		}
	}

	return nil
}
func (m *QuestionModel) InsertInterviewResponse(interviewResponse *QuestionData) error {
	// Validate the question data
	err := m.ValidateQuestionData(interviewResponse)
	if err != nil {
		return err
	}

	// Insert the question data into the database
	query := `INSERT INTO interview_responses (text, type, options, required) VALUES (?, ?, ?, ?)`
	_, err = m.DB.Exec(query, interviewResponse.Text, interviewResponse.Type, interviewResponse.Options, interviewResponse.Required)
	if err != nil {
		return err
	}

	return nil
}
func (m *QuestionModel) GetInterviewResponse(id int) (*QuestionData, error) {
	// Query to get the interview response by ID
	query := `SELECT id, text, type, options, allow_confidence_rating, required FROM interview_responses WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	// Scan the result into a QuestionData struct
	var questionData QuestionData
	err := row.Scan(&questionData.ID, &questionData.Text, &questionData.Type, &questionData.Options, &questionData.Required)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows found
		}
		return nil, err // Other error
	}

	return &questionData, nil
}

func (m *InterviewResponseModel) ValidateInterviewResponse(response *InterviewResponse) error {
	if m == nil {
		return fmt.Errorf("InterviewResponseModel is not initialized")
	}

	if m.Validator == nil {
		return fmt.Errorf("validator is not initialized")
	}

	// Perform validation
	m.Validator.Check(response.QuestionID > 0, "question_id", "Question ID must be a positive integer")
	m.Validator.Check(validator.NotBlank(response.Answer), "answer", "Answer cannot be blank")

	if len(m.Validator.Errors) > 0 {
		return fmt.Errorf("validation errors: %v", m.Validator.Errors)
	}

	return nil
}

func (m *InterviewResponseModel) InsertInterviewResponse(interviewResponse *InterviewResponse) error {
	// Insert the interview response into the database
	query := `INSERT INTO interview_responses (session_id, question_id, response_text, audio_url, confidence, submitted_at) VALUES (?, ?, ?, ?, ?, ?)`
	args := []interface{}{interviewResponse.SessionID, interviewResponse.QuestionID, interviewResponse.Text, interviewResponse.AudioURL, interviewResponse.Confidence, interviewResponse.SubmittedAt}
	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *InterviewResponseModel) GetInterviewResponseByID(id int) (*InterviewResponse, error) {
	// Query to get the interview response by ID
	query := `SELECT id, session_id, question_id, response_text, audio_url, confidence, submitted_at FROM interview_responses WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	// Scan the result into an InterviewResponse struct
	var interviewResponse InterviewResponse
	err := row.Scan(&interviewResponse.ID, &interviewResponse.SessionID, &interviewResponse.QuestionID, &interviewResponse.Text, &interviewResponse.AudioURL, &interviewResponse.Confidence, &interviewResponse.SubmittedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows found
		}
		return nil, err // Other error
	}

	return &interviewResponse, nil
}

// get questionid
func (m *QuestionModel) GetQuestionID(id int) (*QuestionData, error) {
	// Query to get the questionID
	query := `SELECT id, text, type, options, required FROM interview_responses WHERE id = ?`
	row := m.DB.QueryRow(query, id)
	// Scan the result into a QuestionData struct
	var questionData QuestionData
	err := row.Scan(&questionData.ID, &questionData.Text, &questionData.Type, &questionData.Options, &questionData.Required)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows found
		}
		return nil, err // Other error
	}
	return &questionData, nil

}

// get question by id
func (m *QuestionModel) GetQuestionByID(id int) (*QuestionData, error) {
	query := `SELECT id, text, type, options, required FROM questions WHERE id = $1`
	row := m.DB.QueryRow(query, id)

	var question QuestionData
	err := row.Scan(&question.ID, &question.Text, &question.Type, &question.Options, &question.Required)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows found
		}
		return nil, err // Other error
	}

	return &question, nil
}

// get first question ID to start the interview
func (m *QuestionModel) GetFirstQuestionID() (int, error) {
	// Query to get the first question ID
	query := `SELECT id FROM questions ORDER BY id LIMIT 1;`
	row := m.DB.QueryRow(query)

	var questionID int
	err := row.Scan(&questionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No rows found
		}
		return 0, err // Other error
	}

	return questionID, nil
}
func (m *QuestionModel) GetActiveQuestions() ([]*QuestionModel, error) {
	query := `
        SELECT id, text, type, options
        FROM questions
        WHERE is_active = true
        ORDER BY id ASC;
    `

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*QuestionModel
	for rows.Next() {
		var q QuestionModel
		var optionsJSON []byte

		err := rows.Scan(&q.ID, &q.Text, &q.Type)
		if err != nil {
			return nil, err
		}

		// parse options JSON
		if len(optionsJSON) > 0 {
			err = json.Unmarshal(optionsJSON, &q.Options)
			if err != nil {
				return nil, err
			}
		}

		questions = append(questions, &q)
	}

	return questions, nil
}

// save answer function
func (m *InterviewResponseModel) SaveAnswer(interviewResponse *InterviewResponse) error {
	// Validate the interview response
	err := m.ValidateInterviewResponse(interviewResponse)
	if err != nil {
		return err
	}

	// Insert the interview response into the database
	query := `
        INSERT INTO responses (session_id, question_id, response_text, submitted_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err = m.DB.Exec(query, interviewResponse.SessionID, interviewResponse.QuestionID, interviewResponse.Answer, interviewResponse.SubmittedAt)
	return err
}

func (m *InterviewResponseModel) GetInterviewResponsesBySessionID(sessionID int) ([]*InterviewResponse, error) {
	// Query to get interview responses by session ID
	query := `SELECT id, question_id, response_text, audio_url, confidence, submitted_at FROM interview_responses WHERE session_id = ?`
	rows, err := m.DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*InterviewResponse
	for rows.Next() {
		var response InterviewResponse
		err := rows.Scan(&response.ID, &response.QuestionID, &response.Answer)
		if err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

// create a session number for the interview
// set the session number
// set the Teacher ID
// when it started and ended
func (m *InterviewResponseModel) CreateInterviewSession(teacherID int) (int, error) {
	query := `
        INSERT INTO sessions (teacher_id, started_at)
        VALUES ($1, $2)
        RETURNING id
    `
	var sessionID int
	err := m.DB.QueryRow(query, teacherID, time.Now().Format(time.RFC3339)).Scan(&sessionID)
	if err != nil {
		log.Printf("Error creating interview session: %v", err)
		return 0, err
	}
	return sessionID, nil
}

// get all session IDs made by a teacher (user)
func (m *InterviewResponseModel) GetAllSessionsByTeacherID(teacherID int) ([]int, error) {
	query := `SELECT id FROM sessions WHERE teacher_id = $1`
	rows, err := m.DB.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessionIDs []int
	for rows.Next() {
		var sessionID int
		err := rows.Scan(&sessionID)
		if err != nil {
			return nil, err
		}
		sessionIDs = append(sessionIDs, sessionID)
	}

	return sessionIDs, nil
}
