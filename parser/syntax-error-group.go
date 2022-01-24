package parser

type GroupedSyntaxError struct {
	Errors []SyntaxError
}

func NewGroupedSyntaxError(errors []SyntaxError) GroupedSyntaxError {
	return GroupedSyntaxError{
		Errors: errors,
	}
}

func (err GroupedSyntaxError) Error() string {
	errorString := ""
	for _, err := range err.Errors {
		errorString += err.Error()
	}
	return errorString
}
