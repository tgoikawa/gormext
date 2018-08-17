package mysqltype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-sql-driver/mysql"
)

const (
	mySQLNullError uint16 = 1048
)

func assertTimeEquals(t *testing.T, expected, actual time.Time) {
	assert.Truef(t, actual.Equal(expected), "unexpected,  expected: `%s`,actual: `%s`", expected, actual)
}

func assertMySQLErrNumber(t *testing.T, err error, number uint16) {
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			t.Error("err should be MySQL error")
		} else if mysqlError.Number != number {
			t.Errorf("unexpected MySQL error occurred %d", mysqlError.Number)
		}

	} else {
		t.Error("Should be error")
	}

}
