package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// StaffTypeDao - Contact DAO Repository
type StaffTypeDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(staffTypeId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(staffTypeId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(staffTypeId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewStaffTypeDao - Contruct StaffType Dao
func NewStaffTypeDao(client utils.Map, business_id string) StaffTypeDao {
	var daoStaffType StaffTypeDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoStaffType = &mongodb_repository.StaffTypeMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoStaffType != nil {
		// Initialize the Dao
		daoStaffType.InitializeDao(client, business_id)
	}

	return daoStaffType
}
