package mongodb_repository

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-platform/platform_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AttendanceMongoDBDao - Attendance DAO Repository
type AttendanceMongoDBDao struct {
	client     utils.Map
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *AttendanceMongoDBDao) InitializeDao(client utils.Map, bussinesId string, staffId string) {
	log.Println("Initialize Attendance Mongodb DAO")
	p.client = client
	p.businessId = bussinesId
	p.staffId = staffId
}

// ****************************
// List - List all Collections
//
// *****************************

func (p *AttendanceMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map

	log.Println("Begin - Find All Collection Dao", hr_common.DbHrAttendances)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return nil, err
	}

	log.Println("Get Collection - Find All Collection Dao", filter, len(filter), sort, len(sort))

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
	// ==============================================

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

	// Add Lookup stages ================================
	stages = p.appendListLookups(stages)
	// ==================================================

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

// ******************************
// Get - Get designation details
//
// ******************************
func (p *AttendanceMongoDBDao) Get(attendanceId string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("attendanceMongoDao::Get:: Begin ", attendanceId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	log.Println("Find:: Got Collection ")

	filter := bson.D{
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendanceId},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filter = append(filter, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}

	log.Println("Get:: Got filter ", filter)

	singleResult := collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		log.Println("Get:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}
	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("attendanceMongoDao::Get:: End Found a single document: \n", err)
	return result, nil
}

// ********************
// Find - Find by code
//
// ********************
func (p *AttendanceMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("attendanceMongoDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	log.Println("Find:: Got Collection ", err)

	bfilter := bson.D{}
	err = bson.UnmarshalExtJSON([]byte(filter), true, &bfilter)
	if err != nil {
		log.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})
	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		bfilter = append(bfilter, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}

	log.Println("Find:: Got filter ", bfilter)
	singleResult := collection.FindOne(ctx, bfilter)
	if singleResult.Err() != nil {
		log.Println("Find:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("attendanceMongoDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// **************************
// Create - Create Collection
//
// **************************
func (p *AttendanceMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business Attendance Save - Begin", indata)
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return utils.Map{}, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	// Insert a single document
	insertResult, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err
	}

	log.Println("Inserted a single document: ", insertResult.InsertedID)
	log.Println("Save - End", indata[hr_common.FLD_ATTENDANCE_ID])

	//return p.Find(indata[utils.FLD_ATTENDANCE_ID].(string))
	return indata, err
}

// **************************
// Update - Update Collection
//
// **************************
func (p *AttendanceMongoDBDao) Update(attendanceId string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	log.Printf("Update - Values %v", indata)

	filter := bson.D{
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendanceId},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId}}

	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filter = append(filter, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult.ModifiedCount)

	log.Println("Update - End")
	return indata, nil
}

// **************************
// Delete - Delete Collection
//
// **************************
func (p *AttendanceMongoDBDao) Delete(attendanceId string) (int64, error) {

	log.Println("attendanceMongoDao::Delete - Begin ", attendanceId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return 0, err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filter := bson.D{
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendanceId},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId}}

	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filter = append(filter, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}

	res, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("attendanceMongoDao::Delete - End deleted %v documents\n", res.DeletedCount)
	return res.DeletedCount, nil
}
func (p *AttendanceMongoDBDao) appendListLookups(stages []bson.M) []bson.M {

	// Lookup Stage for Token ========================================
	// Lookup Stage
	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":         platform_common.DbPlatformAppUsers,
			"localField":   hr_common.FLD_STAFF_ID,
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
	stages = append(stages, lookupStage)

	// Lookup Stage for User ==========================================
	lookupStage = bson.M{
		"$lookup": bson.M{
			"from":         hr_common.DbHrShifts,
			"localField":   hr_common.FLD_TYPE_OF_WORK,
			"foreignField": hr_common.FLD_SHIFT_ID,
			"as":           hr_common.FLD_SHIFT_INFO,
			"pipeline": []bson.M{
				// Remove following fields from result-set
				{"$project": bson.M{
					db_common.FLD_DEFAULT_ID:              0,
					platform_common.FLD_APP_USER_PASSWORD: 0,
					db_common.FLD_IS_DELETED:              0,
					db_common.FLD_CREATED_AT:              0,
					db_common.FLD_UPDATED_AT:              0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	// Lookup Stage for Token ========================================
	lookupStage = bson.M{
		"$lookup": bson.M{
			"from":         hr_common.DbHrWorkLocations,
			"localField":   hr_common.FLD_WORKLOCATION,
			"foreignField": hr_common.FLD_WORKLOCATION_ID,
			"as":           hr_common.FLD_WORKLOCATION_INFO,
			"pipeline": []bson.M{
				// Remove following fields from result-set
				{"$project": bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	return stages
}
