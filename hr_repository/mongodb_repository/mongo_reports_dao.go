package mongodb_repository

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-platform/platform_common"
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
func (p *ReportsMongoDBDao) GetAttendanceSummary(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

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

	// All Stages
	stages := []bson.M{}

	// Remove unwanted fields =======================
	unsetStage := bson.M{"$unset": db_common.FLD_DEFAULT_ID}
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
	filterdoc = append(filterdoc, bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	matchStage := bson.M{"$match": filterdoc}
	stages = append(stages, matchStage)
	// ==================================================

	// Add Group stage ================================
	groupbyStage := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"staff_id": "$staff_id",
				"for_date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$date_time"}},
			},
			"docs": bson.M{"$push": "$$ROOT"},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, groupbyStage)
	// ==================================================

	// Project Stage =====================================
	projectStage := bson.M{
		"$project": bson.M{
			"docs." + db_common.FLD_CREATED_AT:    0,
			"docs." + db_common.FLD_UPDATED_AT:    0,
			"docs." + db_common.FLD_IS_DELETED:    0,
			"docs." + hr_common.FLD_ATTENDANCE_ID: 0,
			"docs." + hr_common.FLD_BUSINESS_ID:   0,
			"docs." + hr_common.FLD_LATITUDE:      0,
			"docs." + hr_common.FLD_LONGITUDE:     0,
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, projectStage)
	// ==================================================

	// Lookup Stage for staff-info =========================
	lookupStage1 := bson.M{
		"$lookup": bson.M{
			"from":         platform_common.DbPlatformAppUsers,
			"localField":   "_id." + hr_common.FLD_STAFF_ID,
			"foreignField": platform_common.FLD_APP_USER_ID,
			"as":           hr_common.FLD_STAF_INFO,
			"pipeline": []bson.M{
				// Remove following fields from result-set
				{"$project": bson.M{
					db_common.FLD_DEFAULT_ID:              0,
					db_common.FLD_IS_DELETED:              0,
					db_common.FLD_CREATED_AT:              0,
					db_common.FLD_UPDATED_AT:              0,
					platform_common.FLD_APP_USER_PASSWORD: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage1)
	// ==========================================================

	// Lookup Stage for shift =========================
	lookupStage2 := bson.M{
		"$lookup": bson.M{
			"from":         hr_common.DbHrShifts,
			"localField":   "docs.type_of_work",
			"foreignField": hr_common.FLD_SHIFT_ID,
			"as":           hr_common.FLD_SHIFT_INFO,
			"pipeline": []bson.M{
				// Remove following fields from result-set
				{"$project": bson.M{
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
		"$lookup": bson.M{
			"from":         hr_common.DbHrWorkLocations,
			"localField":   "docs.work_location",
			"foreignField": hr_common.FLD_WORKLOCATION_ID,
			"as":           hr_common.FLD_WORKLOCATION_INFO,
			"pipeline": []bson.M{
				// Remove following fields from result-set
				{"$project": bson.M{
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

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("Sort Unmarshal Error ", sort)
		} else {
			sortStage := bson.M{"$sort": sortdoc}
			stages = append(stages, sortStage)
		}
	}

	if skip > 0 {
		skipStage := bson.M{"$skip": skip}
		stages = append(stages, skipStage)
	}

	if limit > 0 {
		limitStage := bson.M{"$limit": limit}
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

	filtercount, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return utils.Map{}, err
	}
	basefilterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		basefilterdoc = append(basefilterdoc, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}
	totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	if err != nil {
		return utils.Map{}, err
	}

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    totalcount,
			db_common.LIST_FILTEREDSIZE: filtercount,
			db_common.LIST_RESULTSIZE:   len(results),
		},
		db_common.LIST_RESULT: results,
	}

	return response, nil
}
