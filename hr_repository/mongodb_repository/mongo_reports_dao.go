package mongodb_repository

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// ReportsMongoDBDao - Reports MongoDB DAO Repository
type ReportsMongoDBDao struct {
	client     utils.Map
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

// InitializeDao - Initialize the ReportsMongoDBDao
func (p *ReportsMongoDBDao) InitializeDao(client utils.Map, businessId string, staffId string) {
	log.Println("Initialize ReportsMongoDBDao")
	p.client = client
	p.businessId = businessId
	p.staffId = staffId
}

// GetAttendanceSummary - Get Attendance Summary data
func (p *ReportsMongoDBDao) GetAttendanceSummary(filter string, aggr string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("Begin - GetAttendanceSummary - Reports - Dao", hr_common.DbHrAttendances)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return nil, err
	}

	log.Println("GetAttendanceSummary - Parameters", filter, len(filter), sort, len(sort))

	filterdoc := bson.D{}
	if len(filter) > 0 {
		// filters, _ := strconv.Unquote(string(filter))
		// 20230803 Karthi: The second parameter should be false to interpret "$date" in JSON
		err = bson.UnmarshalExtJSON([]byte(filter), false, &filterdoc)
		if err != nil {
			log.Println("Unmarshal Ext JSON error", err)
		}
	}

	aggrdoc := bson.D{}
	if len(aggr) > 0 {
		// filters, _ := strconv.Unquote(string(filter))
		// 20230803 Karthi: The second parameter should be false to interpret "$date" in JSON
		err = bson.UnmarshalExtJSON([]byte(aggr), false, &aggrdoc)
		if err != nil {
			log.Println("Unmarshal Ext JSON error", err)
		}
	}

	// All Stages
	stages := []bson.M{}

	// Remove unwanted fields =======================
	unsetStage := bson.M{hr_common.MONGODB_UNSET: db_common.FLD_DEFAULT_ID}
	stages = append(stages, unsetStage)
	// =============================================

	// Match Stage ==================================
	filterdoc = append(filterdoc,
		bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filterdoc = append(filterdoc, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}

	matchStage := bson.M{hr_common.MONGODB_MATCH: filterdoc}
	stages = append(stages, matchStage)
	// ==================================================

	// // Add Group stage ================================
	// groupbyStage := bson.M{
	// 	hr_common.MONGODB_GROUP: bson.M{
	// 		db_common.FLD_DEFAULT_ID: bson.M{
	// 			hr_common.FLD_STAFF_ID: "$" + hr_common.FLD_STAFF_ID,
	// 			"for_date": bson.M{
	// 				hr_common.MONGODB_DATETOSTRING: bson.M{
	// 					hr_common.MONGODB_STR_FORMAT: "%Y-%m-%d", "date": "$" + hr_common.FLD_CLOCK_IN + "." + hr_common.FLD_DATETIME}},
	// 		},
	// 		hr_common.FLD_GROUP_DOCS: bson.M{hr_common.MONGODB_PUSH: "$$ROOT"},
	// 	},
	// }

	if !utils.IsEmpty(aggr) {
		// Add Group stage ================================
		groupbyStage := bson.M{
			hr_common.MONGODB_GROUP: bson.M{
				db_common.FLD_DEFAULT_ID: aggrdoc,
				hr_common.FLD_GROUP_DOCS: bson.M{hr_common.MONGODB_PUSH: hr_common.MONGODB_ROOT},
			},
		}
		// Add it to Aggregate Stage
		stages = append(stages, groupbyStage)
		// ==================================================
	}

	// Project Stage =====================================
	projectStage := bson.M{
		hr_common.MONGODB_PROJECT: bson.M{
			hr_common.FLD_GROUP_DOCS + "." + db_common.FLD_CREATED_AT:  0,
			hr_common.FLD_GROUP_DOCS + "." + db_common.FLD_UPDATED_AT:  0,
			hr_common.FLD_GROUP_DOCS + "." + db_common.FLD_IS_DELETED:  0,
			hr_common.FLD_GROUP_DOCS + "." + hr_common.FLD_BUSINESS_ID: 0,
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, projectStage)
	// ==================================================

	// Add Lookup stages ================================
	stages = p.appendListLookups(stages)
	// ==================================================

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("Sort Unmarshal Error ", sort)
		} else {
			sortStage := bson.M{hr_common.MONGODB_SORT: sortdoc}
			stages = append(stages, sortStage)
		}
	}

	if skip > 0 {
		skipStage := bson.M{hr_common.MONGODB_SKIP: skip}
		stages = append(stages, skipStage)
	}

	if limit > 0 {
		limitStage := bson.M{hr_common.MONGODB_LIMIT: limit}
		stages = append(stages, limitStage)
	}

	cursor, err := collection.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}

	var results []utils.Map
	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// filtercount, err := collection.CountDocuments(ctx, filterdoc)
	// if err != nil {
	// 	return utils.Map{}, err
	// }
	// basefilterdoc := bson.D{
	// 	{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
	// 	{Key: db_common.FLD_IS_DELETED, Value: false}}

	// // Append StaffId in filter if available
	// if len(p.staffId) > 0 {
	// 	basefilterdoc = append(basefilterdoc, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	// }
	// totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	// if err != nil {
	// 	return utils.Map{}, err
	// }

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    len(results), //totalcount,
			db_common.LIST_FILTEREDSIZE: len(results), //filtercount,
			db_common.LIST_RESULTSIZE:   len(results),
		},
		db_common.LIST_RESULT: results,
	}

	return response, nil
}

func (p *ReportsMongoDBDao) appendListLookups(stages []bson.M) []bson.M {

	// // Lookup Stage for staff-info =========================
	// lookupStage1 := bson.M{
	// 	hr_common.MONGODB_LOOKUP: bson.M{
	// 		hr_common.MONGODB_STR_FROM:         platform_common.DbPlatformAppUsers,
	// 		hr_common.MONGODB_STR_LOCALFIELD:   "_id." + hr_common.FLD_STAFF_ID,
	// 		hr_common.MONGODB_STR_FOREIGNFIELD: platform_common.FLD_APP_USER_ID,
	// 		hr_common.MONGODB_STR_AS:           hr_common.FLD_STAFF_INFO,
	// 		hr_common.MONGODB_STR_PIPELINE: []bson.M{
	// 			// Remove following fields from result-set
	// 			{hr_common.MONGODB_PROJECT: bson.M{
	// 				db_common.FLD_DEFAULT_ID:              0,
	// 				db_common.FLD_IS_DELETED:              0,
	// 				db_common.FLD_CREATED_AT:              0,
	// 				db_common.FLD_UPDATED_AT:              0,
	// 				platform_common.FLD_APP_USER_PASSWORD: 0}},
	// 		},
	// 	},
	// }
	// // Add it to Aggregate Stage
	// stages = append(stages, lookupStage1)
	// // ==========================================================

	// Lookup Stage for shift =========================
	lookupStage2 := bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         hr_common.DbHrShiftProfiles,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_GROUP_DOCS + "." + hr_common.FLD_CLOCK_IN + "." + "type_of_work",
			hr_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_SHIFT_PROFILE_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_SHIFT_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID:  0,
					db_common.FLD_IS_DELETED:  0,
					db_common.FLD_CREATED_AT:  0,
					db_common.FLD_UPDATED_AT:  0,
					hr_common.FLD_BUSINESS_ID: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage2)
	// ==========================================================

	// Lookup Stage for Work Location =========================
	lookupStage3 := bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         hr_common.DbHrWorkLocations,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_GROUP_DOCS + "." + hr_common.FLD_CLOCK_IN + "." + "work_location",
			hr_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_WORKLOCATION_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_WORKLOCATION_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID:  0,
					db_common.FLD_IS_DELETED:  0,
					db_common.FLD_CREATED_AT:  0,
					db_common.FLD_UPDATED_AT:  0,
					hr_common.FLD_BUSINESS_ID: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage3)
	// ==========================================================

	// Lookup Stage for project =========================
	lookupStage4 := bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         hr_common.DbHrProjects,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_GROUP_DOCS + "." + hr_common.FLD_CLOCK_IN + "." + hr_common.FLD_PROJECT_ID,
			hr_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_PROJECT_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_PROJECT_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID:  0,
					db_common.FLD_IS_DELETED:  0,
					db_common.FLD_CREATED_AT:  0,
					db_common.FLD_UPDATED_AT:  0,
					hr_common.FLD_BUSINESS_ID: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage4)
	// ==========================================================

	// Lookup Stage for client =========================
	lookupStage5 := bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         hr_common.DbHrClients,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_GROUP_DOCS + "." + hr_common.FLD_CLOCK_IN + "." + hr_common.FLD_CLIENT_ID,
			hr_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_CLIENT_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_CLIENT_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID:  0,
					db_common.FLD_IS_DELETED:  0,
					db_common.FLD_CREATED_AT:  0,
					db_common.FLD_UPDATED_AT:  0,
					hr_common.FLD_BUSINESS_ID: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage5)
	// ==========================================================

	return stages
}
