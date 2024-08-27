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
	// api
	SuccessToCreate = "crud.successToCreate"
	SuccessToRead   = "crud.successToRead"
	SuccessToUpdate = "crud.successToUpdate"
	SuccessToDelete = "crud.successToDelete"
	// db
	DbModelCreatedError = "error.db.model.create"
	DbModelReadError    = "error.db.model.read"
	DbModelUpdatedError = "error.db.model.update"
	DbModelDeletedError = "error.db.model.delete"
	// update related
	DbModelUpdatedIdCannotBeZero  = "error.db.model.update.idCannotBeZero"
	DbModelAlreadyUpdatedByOthers = "error.db.model.update.alreadyUpdatedByOthers"
	// read related
	DbModelReadNotExists     = "error.db.model.read.notExists"     // {{.Id}}
	DbModelReadFieldNotFound = "error.db.model.read.fieldNotFound" // {{.Field}}
	DbModelReadOprNotFound   = "error.db.model.read.oprNotFound"   // {{.Opr}}
)

// bind
const (
	MissingRequiredParam = "bing.url.missingRequiredParam"
	InvalidUintParam     = "bing.url.invalidUintParam"
)
