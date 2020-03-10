package errorp

type Error interface {
	Error() string
	Code() string
	Status() int64
	Description() string
	Hint() string
	WithCode(c string) Error
	WithDescription(d string) Error
	WithHint(h string) Error
}

type PolicyError struct {
	C string `json:"error_code"`
	D string `json:"error_description"`
	H string `json:"error_hint"`
	S int64  `json:"-"`
}

func NewPolicyError(statusCode int64, code, hint, description string) *PolicyError {
	e := PolicyError{S: statusCode, H: hint, D: description, C: code}
	return &e
}

func (e *PolicyError) Error() string {
	return "Error " + e.C + ": " + e.D
}

func (e *PolicyError) Code() string {
	return e.C
}
func (e *PolicyError) Description() string {
	return e.D
}
func (e *PolicyError) Hint() string {
	return e.H
}
func (e *PolicyError) WithCode(c string) Error {
	e.C = c
	return e
}
func (e *PolicyError) WithDescription(d string) Error {
	e.D = d
	return e
}
func (e *PolicyError) WithHint(h string) Error {
	e.H = h
	return e
}
func (e *PolicyError) Status() int64 {
	return e.S
}
