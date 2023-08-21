package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// DasboardDao - User DAO Repository
type DashboardDao interface {
	InitializeDao(client utils.Map, businessId string, staffId string)
	GetDashboardData() (utils.Map, error)
}

func NewDashboardDao(client utils.Map, businessId string, staffId string) DashboardDao {
	var daoDashboard DashboardDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbDashboard, _ := db_common.GetDatabaseType(client)

	switch dbDashboard {
	case db_common.DATABASE_TYPE_MONGODB:
		daoDashboard = &mongodb_repository.DashboardMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
		daoDashboard = nil
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
		daoDashboard = nil
	}

	if daoDashboard != nil {
		// Initialize the Dao
		daoDashboard.InitializeDao(client, businessId, staffId)
	}

	return daoDashboard
}
