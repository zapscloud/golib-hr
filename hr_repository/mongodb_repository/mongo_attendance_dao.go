package mongodb_repository

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AttendanceMongoDBDao - Attendance DAO Repository
type AttendanceMongoDBDao struct {
	client     utils.Map
	businessID string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *AttendanceMongoDBDao) InitializeDao(client utils.Map, bussinesId string) {
	log.Println("Initialize Attendance Mongodb DAO")
	p.client = client
	p.businessID = bussinesId
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

	opts := options.Find()

	filterdoc := bson.D{}
	if len(filter) > 0 {
		// filters, _ := strconv.Unquote(string(filter))
		err = bson.UnmarshalExtJSON([]byte(filter), true, &filterdoc)
		if err != nil {
			log.Println("Unmarshal Ext JSON error", err)
			log.Println(filterdoc)
		}
	}

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("Sort Unmarshal Error ", sort)
		} else {
			opts.SetSort(sortdoc)
		}
	}

	if skip > 0 {
		log.Println(filterdoc)
		opts.SetSkip(skip)
	}

	if limit > 0 {
		log.Println(filterdoc)
		opts.SetLimit(limit)
	}

	filterdoc = append(filterdoc,
		bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})
	log.Println("Parameter values ", filterdoc, opts)
	cursor, err := collection.Find(ctx, filterdoc, opts)
	if err != nil {
		return nil, err
	}

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	log.Println("End - Find All Collection Dao", results)

	listdata := []utils.Map{}
	for idx, value := range results {
		log.Println("Item ", idx)
		value = db_common.AmendFldsForGet(value)
		listdata = append(listdata, value)
	}

	log.Println("End - Find All Collection Dao", listdata)

	log.Println("Parameter values ", filterdoc)
	filtercount, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return nil, err
	}

	basefilterdoc := bson.D{{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID}}
	totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	if err != nil {
		return nil, err
	}

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    totalcount,
			db_common.LIST_FILTEREDSIZE: filtercount,
			db_common.LIST_RESULTSIZE:   len(listdata),
		},
		db_common.LIST_RESULT: listdata,
	}

	return response, nil

}

// ******************************
// Get - Get designation details
//
// ******************************
func (p *AttendanceMongoDBDao) Get(attendance_id string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("attendanceMongoDao::Get:: Begin ", attendance_id)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	log.Println("Find:: Got Collection ")

	filter := bson.D{
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendance_id},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

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
		fmt.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

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
func (p *AttendanceMongoDBDao) Update(attendance_id string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrAttendances)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	log.Printf("Update - Values %v", indata)

	filter := bson.D{
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendance_id},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

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
func (p *AttendanceMongoDBDao) Delete(attendance_id string) (int64, error) {

	log.Println("attendanceMongoDao::Delete - Begin ", attendance_id)

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
		{Key: hr_common.FLD_ATTENDANCE_ID, Value: attendance_id},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessID}}

	res, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("attendanceMongoDao::Delete - End deleted %v documents\n", res.DeletedCount)
	return res.DeletedCount, nil
}
