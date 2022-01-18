package sql

import (
	"strings"
	"swallowtail/libraries/gerrors"
)

// PostgresGetFailed ...
func PostgresGetFailed(err error) error {
	return gerrors.Augment(gerrors.New(gerrors.ErrUnknown, strings.ReplaceAll(err.Error(), " ", "_"), nil), "postgres_get_failed", nil)
}

// PostgresSelectFailed ...
func PostgresSelectFailed(err error) error {
	return gerrors.Augment(gerrors.New(gerrors.ErrUnknown, strings.ReplaceAll(err.Error(), " ", "_"), nil), "postgres_select_failed", nil)
}

// PostgresExecFailed ..
func PostgresExecFailed(err error) error {
	return gerrors.Augment(gerrors.New(gerrors.ErrUnknown, strings.ReplaceAll(err.Error(), " ", "_"), nil), "postgres_exec_failed", nil)
}

// PostgresNoRowsFound ...
func PostgresNoRowsFound(err error) error {
	return gerrors.Augment(gerrors.New(gerrors.ErrNotFound, strings.ReplaceAll(err.Error(), " ", "_"), nil), "postgres_no_rows_found", nil)
}

// PostgresManyRowsFound ...
func PostgresManyRowsFound(err error) error {
	return gerrors.Augment(gerrors.New(gerrors.ErrAlreadyExists, strings.ReplaceAll(err.Error(), " ", "_"), nil), "postgres_many_rows_found", nil)
}
