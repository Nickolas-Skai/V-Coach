package main

import (

)

// IsValidQuestionType checks if the provided question type is valid.
func IsValidQuesgtionType(questionType string) bool {
    validTypes := []string{"text", "checkbox", "radio", "scale"}
    for _, validType := range validTypes {
        if questionType == validType {
            return true
        }
    }
    return false
}