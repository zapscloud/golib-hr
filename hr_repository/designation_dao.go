package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// DesignationDao - Designation DAO Repository
type DesignationDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Designation Details
	Get(designation_id string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Designation
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(designation_id string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(designation_id string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)
}

// NewdesignationMongoDao - Contruct Designation Dao
func NewDesignationDao(client utils.Map, businessid string) DesignationDao {
	var daoDesignation DesignationDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoDesignation = &mongodb_repository.DesignationMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoDesignation != nil {
		// Initialize the Dao
		daoDesignation.InitializeDao(client, businessid)
	}

	return daoDesignation
}
