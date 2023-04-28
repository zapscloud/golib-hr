package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// WorkLocationDao - Contact DAO Repository
type WorkLocationDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(workLocId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(workLocId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(workLocId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewWorkLocationDao - Contruct Holiday Dao
func NewWorkLocationDao(client utils.Map, business_id string) WorkLocationDao {
	var daoWorkLocation WorkLocationDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoWorkLocation = &mongodb_repository.WorkLocationMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoWorkLocation != nil {
		// Initialize the Dao
		daoWorkLocation.InitializeDao(client, business_id)
	}

	return daoWorkLocation
}
