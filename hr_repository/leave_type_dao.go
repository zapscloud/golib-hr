package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// LeaveTypeDao - LeaveType DAO Repository
type LeaveTypeDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get LeaveType Details
	Get(LeaveTypeid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create LeaveType
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(LeaveTypeid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(LeaveTypeid string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)

	GetDeptCodeDetails(LeaveTypecode string) (utils.Map, error)
}

// NewLeaveTypeDao - Contruct LeaveType Dao
func NewLeaveTypeDao(client utils.Map, businessid string) LeaveTypeDao {
	var daoLeaveType LeaveTypeDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoLeaveType = &mongodb_repository.LeaveTypeMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoLeaveType != nil {
		// Initialize the Dao
		daoLeaveType.InitializeDao(client, businessid)
	}

	return daoLeaveType
}
