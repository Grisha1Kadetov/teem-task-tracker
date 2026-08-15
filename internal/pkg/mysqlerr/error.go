package mysqlerr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

const duplicateEntryNumber uint16 = 1062

var duplicateEntry = &mysql.MySQLError{Number: duplicateEntryNumber}

func IsDuplicateEntry(err error) bool {
	return errors.Is(err, duplicateEntry)
}
