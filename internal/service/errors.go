package service

type OperationErrors struct {
	Validation map[string][]string
	Database   map[string][]string
}

func (ve *OperationErrors) AddValidationError(field, message string) {
	if ve.Validation == nil {
		ve.Validation = make(map[string][]string)
	}
	ve.Validation[field] = append(ve.Validation[field], message)
}

func (ve *OperationErrors) AddDatabaseError(field, message string) {
	if ve.Database == nil {
		ve.Database = make(map[string][]string)
	}
	ve.Database[field] = append(ve.Database[field], message)
}
