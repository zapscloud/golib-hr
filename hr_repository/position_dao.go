package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// PositionDao - Contact DAO Repository
type PositionDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(position_id string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(position_id string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(position_id string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewPositionMongoDao - Contruct Position Dao
func NewPositionDao(client utils.Map, business_id string) PositionDao {
	var daoPosition PositionDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoPosition = &mongodb_repository.PositionMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoPosition != nil {
		// Initialize the Dao
		daoPosition.InitializeDao(client, business_id)
	}

	return daoPosition
}
