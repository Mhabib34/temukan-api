package exception

// File ini berisi semua custom error type yang dipakai di project.
// Tambahkan di sini jika belum ada di exception package yang sudah ada.

// ── Error Types ───────────────────────────────────────────────────────────────

type NotFoundError struct{ Message string }
type UnauthorizedError struct{ Message string }
type ForbiddenError struct{ Message string }
type ConflictError struct{ Message string }
type BadRequestError struct{ Message string }

// Gunakan value receiver agar type assertion err.(ForbiddenError) di ErrorHandler bisa match
func (e NotFoundError) Error() string     { return e.Message }
func (e UnauthorizedError) Error() string { return e.Message }
func (e ForbiddenError) Error() string    { return e.Message }
func (e ConflictError) Error() string     { return e.Message }
func (e BadRequestError) Error() string   { return e.Message }

func NewNotFoundError(msg string) NotFoundError         { return NotFoundError{msg} }
func NewUnauthorizedError(msg string) UnauthorizedError { return UnauthorizedError{msg} }
func NewForbiddenError(msg string) ForbiddenError       { return ForbiddenError{msg} }
func NewConflictError(msg string) ConflictError         { return ConflictError{msg} }
func NewBadRequestError(msg string) BadRequestError     { return BadRequestError{msg} }

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
