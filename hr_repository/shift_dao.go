package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// ShiftDao - Contact DAO Repository
type ShiftDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(shiftId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(shiftId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(shiftId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewShiftDao - Contruct Leave Dao
func NewShiftDao(client utils.Map, businessId string) ShiftDao {
	var daoShift ShiftDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoShift = &mongodb_repository.ShiftMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoShift != nil {
		// Initialize the Dao
		daoShift.InitializeDao(client, businessId)
	}

	return daoShift
}
