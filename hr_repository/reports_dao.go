package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// ReportsDao - Reports DAO Repository
type ReportsDao interface {
	InitializeDao(client utils.Map, businessId string, staffId string)
	GetAttendanceSummary(filter string, sort string, skip int64, limit int64) (utils.Map, error)
}

func NewReportsDao(client utils.Map, businessId string, staffId string) ReportsDao {
	var daoReports ReportsDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbReports, _ := db_common.GetDatabaseType(client)

	switch dbReports {
	case db_common.DATABASE_TYPE_MONGODB:
		daoReports = &mongodb_repository.ReportsMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
		daoReports = nil
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
		daoReports = nil
	}

	if daoReports != nil {
		// Initialize the Dao
		daoReports.InitializeDao(client, businessId, staffId)
	}

	return daoReports
}
