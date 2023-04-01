package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// DepartmentDao - Department DAO Repository
type DepartmentDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Department Details
	Get(departmentid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Department
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(departmentid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(departmentid string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)

	GetDeptCodeDetails(departmentcode string) (utils.Map, error)
}

// NewDepartmentDao - Contruct Department Dao
func NewDepartmentDao(client utils.Map, businessid string) DepartmentDao {
	var daoDepartment DepartmentDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoDepartment = &mongodb_repository.DepartmentMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoDepartment != nil {
		// Initialize the Dao
		daoDepartment.InitializeDao(client, businessid)
	}

	return daoDepartment
}
