package mongodb_repository

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// DashboardMongoDBDao - Dashboard MongoDB DAO Repository
type DashboardMongoDBDao struct {
	client     utils.Map
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

// InitializeDao - Initialize the DashboardMongoDBDao
func (p *DashboardMongoDBDao) InitializeDao(client utils.Map, businessId string, staffId string) {
	log.Println("Initialize DashboardMongoDBDao")
	p.client = client
	p.businessId = businessId
	p.staffId = staffId
}

// GetDashboardData - Get dashboard data
func (p *DashboardMongoDBDao) GetDashboardData() (utils.Map, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return nil, err
	}
	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filterdoc = append(filterdoc, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}
	// 1. Find Total number of Tokens
	totalcount, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return nil, err
	}

	// 2. Count different leave types
	leaveCounts := make(map[string]int64)

	leaveTypes := []string{"sick Leave", "Casual Leave", "Permission", "Leave"}
	for _, leaveType := range leaveTypes {

		// Append each leave_types in filterdoc
		leaveTypefilterDoc := append(filterdoc, bson.E{Key: "leave_type", Value: leaveType})
		count, err := collection.CountDocuments(ctx, leaveTypefilterDoc)
		if err != nil {
			return nil, err
		}
		leaveCounts[leaveType] = count
	}

	// Prepare return data
	retData := utils.Map{
		"Leave":        leaveCounts["Leave"],
		"total_leave":  totalcount,
		"Permission":   leaveCounts["Permission"],
		"sickleave":    leaveCounts["sick Leave"],
		"Casual Leave": leaveCounts["Casual Leave"],
	}

	return retData, nil
}
