package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in            interface{}
		expectedErr   error
		validationErr ValidationErrors
	}{
		{
			in: User{
				ID:     "100500",
				Name:   "Вася",
				Age:    60,
				Email:  "@bk.ru",
				Role:   "admin",
				Phones: []string{"790000000009", "79001111111", "79002222222"},
				meta:   json.RawMessage("{}"),
			},
			expectedErr: nil,
			validationErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidationLengthValue,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidationMoreThanMaximalValue,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidationRegexpValue,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidationLengthValue,
				},
			},
		},
		{
			in: App{
				Version: "1.1.0",
			},
			expectedErr:   nil,
			validationErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 201,
				Body: "{}",
			},
			expectedErr: nil,
			validationErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidationNotInAgreeValues,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			validationErr, expectedErr := Validate(tt.in)
			require.Equal(t, tt.validationErr, validationErr)
			require.Equal(t, tt.expectedErr, expectedErr)
			_ = tt
		})
	}
}

func TestNotStruct(t *testing.T) {
	tests := []struct {
		in            interface{}
		expectedErr   error
		validationErr ValidationErrors
	}{
		{
			in:            1,
			expectedErr:   ErrFormatData,
			validationErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			validationErr, expectedErr := Validate(tt.in)
			require.Equal(t, tt.validationErr, validationErr)
			require.Equal(t, tt.expectedErr, expectedErr)
			_ = tt
		})
	}
}
