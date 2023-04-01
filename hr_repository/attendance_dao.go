package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// AttendanceDao - Attendance DAO Repository
type AttendanceDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Attendance Details
	Get(attendance_id string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Attendance
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(attendance_id string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(attendance_id string) (int64, error)
}

// NewattendanceMongoDao - Contruct Attendance Dao
func NewAttendanceDao(client utils.Map, businessid string) AttendanceDao {
	var daoAttendance AttendanceDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoAttendance = &mongodb_repository.AttendanceMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoAttendance != nil {
		// Initialize the Dao
		daoAttendance.InitializeDao(client, businessid)
	}

	return daoAttendance
}
