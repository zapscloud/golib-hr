package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// LeaveDao - Contact DAO Repository
type LeaveDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string, staffId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(leaveId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(leaveId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(leaveId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewLeaveDao - Contruct Leave Dao
func NewLeaveDao(client utils.Map, businessId string, staffId string) LeaveDao {
	var daoLeave LeaveDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoLeave = &mongodb_repository.LeaveMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoLeave != nil {
		// Initialize the Dao
		daoLeave.InitializeDao(client, businessId, staffId)
	}

	return daoLeave
}
