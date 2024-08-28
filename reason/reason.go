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
	SuccessToCreate   = "base.success.api.create"
	SuccessToRead     = "base.success.api.read"
	SuccessToUpdate   = "base.success.api.update"
	SuccessToDelete   = "base.success.api.delete"
	SuccessToSave     = "base.success.api.save"
	SuccessToSend     = "base.success.api.send"
	SuccessToRegister = "base.success.api.register"
	SuccessToLogin    = "base.success.api.login"
	// db
	DbModelCreatedError = "error.db.model.create.common"
	DbModelReadError    = "error.db.model.read.common"
	DbModelUpdatedError = "error.db.model.update.common"
	DbModelDeletedError = "error.db.model.delete.common"
	// create related
	DbModelCreateDupKey = "error.db.model.create.dupKey" // {{.Key}} uniqueIndex的列有重复值
	// update related
	DbModelUpdatedIdCannotBeZero  = "error.db.model.update.idCannotBeZero"
	DbModelAlreadyUpdatedByOthers = "error.db.model.update.alreadyUpdatedByOthers"
	// read related
	DbModelReadIdNotExists       = "error.db.model.read.idNotExists"       // {{.Id}}
	DbModelReadStrFieldNotExists = "error.db.model.read.strFieldNotExists" // {{.StrField}}
)

// bind
const (
	MissingRequiredParam = "error.bind.url.missingRequiredParam" // {{.Param}}
	InvalidUintParam     = "error.bind.url.invalidUintParam"     // {{.Param}}
)
