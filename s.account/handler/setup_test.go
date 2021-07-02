package handler

import (
	"swallowtail/s.account/dao"
	"testing"
)

func TestMain(m *testing.M) {
	dao.WithMock()
}
