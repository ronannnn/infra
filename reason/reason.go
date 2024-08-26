package reason

const (
	// Success .
	Success = "base.success"
	// UnknownError unknown error
	UnknownError = "base.unknown"
	// RequestFormatError request format error
	RequestFormatError = "base.requestFormatError"
	// UnauthorizedError unauthorized error
	UnauthorizedError = "base.unauthorizedError"
	// DatabaseError database error
	DatabaseError = "base.databaseError"
	// ForbiddenError forbidden error
	ForbiddenError = "base.forbiddenError"
	// DuplicateRequestError duplicate request error
	DuplicateRequestError = "base.duplicateRequestError"
)

// internal error
const (
	ValidatorLangNotFound = "error.i18n.validatorLangNotFound"
)

// CRUD
const (
	SuccessToCreate = "crud.successToCreate"
	SuccessToRead   = "crud.successToRead"
	SuccessToUpdate = "crud.successToUpdate"
	SuccessToDelete = "crud.successToDelete"
)

// bind
const (
	MissingRequiredParam = "bing.url.missingRequiredParam"
	InvalidUintParam     = "bing.url.invalidUintParam"
)
