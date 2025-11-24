package validator

import (
	"testing"
)

func TestPasswordValidator_Validate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid strong password",
			password: "MySecureP@ssw0rd",
			wantErr:  nil,
		},
		{
			name:     "valid minimum length",
			password: "abcd1234",
			wantErr:  nil,
		},
		{
			name:     "valid long password",
			password: "ThisIsAVeryLongPasswordThatIsStillValidAndSecure123456",
			wantErr:  nil,
		},
		{
			name:     "valid with special characters",
			password: "P@ssw0rd!#$%",
			wantErr:  nil,
		},
		{
			name:     "valid with spaces in middle",
			password: "my secure password 123",
			wantErr:  nil,
		},

		{
			name:     "too short - 7 characters",
			password: "pass123",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "too short - empty",
			password: "",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "too short - 1 character",
			password: "a",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "too long - 129 characters",
			password: "a" + string(make([]byte, 128)),
			wantErr:  ErrPasswordTooLong,
		},

		{
			name:     "common password - password",
			password: "password",
			wantErr:  ErrPasswordTooCommon,
		},
		{
			name:     "common password - 12345678",
			password: "12345678",
			wantErr:  ErrPasswordTooCommon,
		},
		{
			name:     "common password - password123",
			password: "password123",
			wantErr:  ErrPasswordTooCommon,
		},
		{
			name:     "common password - iloveyou",
			password: "iloveyou",
			wantErr:  ErrPasswordTooCommon,
		},
		{
			name:     "common password - uppercase PASSWORD",
			password: "PASSWORD",
			wantErr:  ErrPasswordTooCommon,
		},
		{
			name:     "common password - mixed case PaSsWoRd",
			password: "PaSsWoRd",
			wantErr:  ErrPasswordTooCommon,
		},

		{
			name:     "all numeric - 8 digits",
			password: "87654321",
			wantErr:  ErrPasswordAllNumeric,
		},
		{
			name:     "all numeric - long",
			password: "123456789012345",
			wantErr:  ErrPasswordAllNumeric,
		},

		{
			name:     "excessive repeating - aaaa",
			password: "aaaaaaaa",
			wantErr:  ErrPasswordRepeating,
		},
		{
			name:     "excessive repeating - in middle",
			password: "pass1111word",
			wantErr:  ErrPasswordRepeating,
		},
		{
			name:     "excessive repeating - at end",
			password: "password1111",
			wantErr:  ErrPasswordRepeating,
		},
		{
			name:     "acceptable repeating - 3 consecutive",
			password: "passs123",
			wantErr:  nil,
		},
		{
			name:     "acceptable repeating - non-consecutive",
			password: "abababab",
			wantErr:  nil,
		},

		{
			name:     "leading whitespace",
			password: " password123",
			wantErr:  ErrPasswordWhitespace,
		},
		{
			name:     "trailing whitespace",
			password: "password123 ",
			wantErr:  ErrPasswordWhitespace,
		},
		{
			name:     "both leading and trailing whitespace",
			password: " password123 ",
			wantErr:  ErrPasswordWhitespace,
		},
		{
			name:     "multiple leading spaces",
			password: "   password123",
			wantErr:  ErrPasswordWhitespace,
		},
		{
			name:     "multiple trailing spaces",
			password: "password123   ",
			wantErr:  ErrPasswordWhitespace,
		},

		{
			name:     "exactly 8 characters",
			password: "abcdefgh",
			wantErr:  nil,
		},
		{
			name:     "exactly 128 characters",
			password: "alkfdjalskdfjal;skdjfa;kljf alksjfda;lkfq321jfal;kdjc alsdfj;lkafdsjdlkaj fslkdjasldkfja;lskdfj;alksjdf;alksjdf;alksjdf;alksjdwd",
			wantErr:  nil,
		},
		{
			name:     "unicode characters",
			password: "pässwörd123",
			wantErr:  nil,
		},
		{
			name:     "emoji in password",
			password: "password😀123",
			wantErr:  nil,
		},
	}

	validator := NewPasswordValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)

			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "SecurePass123",
			wantErr:  nil,
		},
		{
			name:     "invalid - too short",
			password: "short1",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "invalid - common",
			password: "password",
			wantErr:  ErrPasswordTooCommon,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)

			if err != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordValidator_CustomConfig(t *testing.T) {
	tests := []struct {
		name      string
		validator *PasswordValidator
		password  string
		wantErr   error
	}{
		{
			name: "custom min length - 12",
			validator: &PasswordValidator{
				MinLength:       12,
				MaxLength:       128,
				CheckCommon:     false,
				CheckRepeating:  false,
				CheckAllNumeric: false,
				CheckWhitespace: false,
			},
			password: "short123456",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name: "disable common check",
			validator: &PasswordValidator{
				MinLength:       8,
				MaxLength:       128,
				CheckCommon:     false,
				CheckRepeating:  false,
				CheckAllNumeric: false,
				CheckWhitespace: false,
			},
			password: "password",
			wantErr:  nil,
		},
		{
			name: "disable repeating check",
			validator: &PasswordValidator{
				MinLength:       8,
				MaxLength:       128,
				CheckCommon:     false,
				CheckRepeating:  false,
				CheckAllNumeric: false,
				CheckWhitespace: false,
			},
			password: "aaaaaaaaaa",
			wantErr:  nil,
		},
		{
			name: "disable all numeric check",
			validator: &PasswordValidator{
				MinLength:       8,
				MaxLength:       128,
				CheckCommon:     false,
				CheckRepeating:  false,
				CheckAllNumeric: false,
				CheckWhitespace: false,
			},
			password: "12345678",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate(tt.password)

			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordValidator_isCommonPassword(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"common - password", "password", true},
		{"common - uppercase", "PASSWORD", true},
		{"common - mixed case", "PaSsWoRd", true},
		{"common - qwerty", "qwerty", true},
		{"not common", "UniqueP@ssw0rd", false},
		{"not common - similar to common", "password1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.isCommonPassword(tt.password)
			if got != tt.want {
				t.Errorf("isCommonPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordValidator_isAllNumeric(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"all numeric", "12345678", true},
		{"all numeric - long", "123456789012345", true},
		{"not numeric - letters", "abc12345", false},
		{"not numeric - special chars", "123!456", false},
		{"empty", "", false},
		{"single digit", "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.isAllNumeric(tt.password)
			if got != tt.want {
				t.Errorf("isAllNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordValidator_hasExcessiveRepeating(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"4 consecutive same char", "aaaa", true},
		{"5 consecutive same char", "aaaaa", true},
		{"4 consecutive in middle", "pass1111word", true},
		{"3 consecutive - acceptable", "aaa", false},
		{"3 consecutive in password", "pass111word", false},
		{"non-consecutive repeating", "abababab", false},
		{"short password", "abc", false},
		{"no repeating", "abcdefgh", false},
		{"4 at end", "password1111", true},
		{"4 at start", "1111password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.hasExcessiveRepeating(tt.password)
			if got != tt.want {
				t.Errorf("hasExcessiveRepeating() = %v, want %v", got, tt.want)
			}
		})
	}
}
