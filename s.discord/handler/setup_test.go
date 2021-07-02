package handler

import (
	"swallowtail/s.discord/dao"
	"testing"
)

func TestMain(m *testing.M) {
	dao.WithMock()
}
