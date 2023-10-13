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

// LeaveMongoDBDao - Leave DAO Repository
type LeaveMongoDBDao struct {
	client     utils.Map
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *LeaveMongoDBDao) InitializeDao(client utils.Map, businessId string, staffId string) {
	log.Println("Initialize Leave Mongodb DAO")
	p.client = client
	p.businessId = businessId
	p.staffId = staffId
}

// List - List all Collections
func (p *LeaveMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map
	var bFilter bool = false
	log.Println("Begin - Find All Collection Dao", hr_common.DbHrLeaves)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
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
		bFilter = true
	}

	// All Stages
	stages := []bson.M{}

	// Remove unwanted fields
	unsetStage := bson.M{hr_common.MONGODB_UNSET: db_common.FLD_DEFAULT_ID}
	stages = append(stages, unsetStage)

	// Add Lookup stages
	stages = p.appendListLookups(stages)

	// Match Stage
	filterdoc = append(filterdoc,
		bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})
	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filterdoc = append(filterdoc, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}
	filterdoc = append(filterdoc, bson.E{Key: db_common.FLD_IS_DELETED, Value: false})
	matchStage := bson.M{hr_common.MONGODB_MATCH: filterdoc}
	stages = append(stages, matchStage)

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

	var filtercount int64 = 0
	if bFilter {
		// Prepare Filter Stages
		filterStages := stages

		// Add Count aggregate
		countStage := bson.M{hr_common.MONGODB_COUNT: hr_common.FLD_FILTERED_COUNT}
		filterStages = append(filterStages, countStage)

		//log.Println("Aggregate for Count ====>", filterStages, stages)

		// Execute aggregate to find the count of filtered_size
		cursor, err := collection.Aggregate(ctx, filterStages)
		if err != nil {
			log.Println("Error in Aggregate", err)
			return nil, err
		}
		var countResult []utils.Map
		if err = cursor.All(ctx, &countResult); err != nil {
			log.Println("Error in cursor.all", err)
			return nil, err
		}

		//log.Println("Count Filter ===>", countResult)
		if len(countResult) > 0 {
			if dataVal, dataOk := countResult[0][hr_common.FLD_FILTERED_COUNT]; dataOk {
				filtercount = int64(dataVal.(int32))
			}
		}
		// log.Println("Count ===>", filtercount)

	} else {
		filtercount, err = collection.CountDocuments(ctx, filterdoc)
		if err != nil {
			return nil, err
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

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
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

	if results == nil {
		results = []utils.Map{}
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

// Get - Get account details
func (p *LeaveMongoDBDao) Get(account_id string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("accountMongoDao::Get:: Begin ", account_id)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	log.Println("Find:: Got Collection ")

	filter := bson.D{
		{Key: hr_common.FLD_LEAVE_ID, Value: account_id},
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

	log.Println("accountMongoDao::Get:: End Found a single document: \n", err)
	return result, nil
}

// Find - Find by code
func (p *LeaveMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("accountMongoDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
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

	log.Println("accountMongoDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (p *LeaveMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business Leave Save - Begin", indata)
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return indata, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	// Insert a single document
	insertResult, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert ", err)
		return indata, err

	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	log.Println("Save - End", indata[hr_common.FLD_LEAVE_ID])

	return indata, err
}

// Update - Update Collection
func (p *LeaveMongoDBDao) Update(account_id string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	log.Printf("Update - Values %v", indata)

	filter := bson.D{
		{Key: hr_common.FLD_LEAVE_ID, Value: account_id},
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId}}
	// Append StaffId in filter if available
	if len(p.staffId) > 0 {
		filter = append(filter, bson.E{Key: hr_common.FLD_STAFF_ID, Value: p.staffId})
	}
	updateResult, err := collection.UpdateOne(ctx, filter, bson.D{{Key: hr_common.MONGODB_SET, Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult.ModifiedCount)

	log.Println("Update - End")
	return indata, nil
}

// Delete - Delete Collection
func (p *LeaveMongoDBDao) Delete(account_id string) (int64, error) {

	log.Println("accountMongoDao::Delete - Begin ", account_id)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return 0, err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filter := bson.D{
		{Key: hr_common.FLD_LEAVE_ID, Value: account_id},
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
	log.Printf("accountMongoDao::Delete - End deleted %v documents\n", res.DeletedCount)
	return res.DeletedCount, nil
}

// DeleteAll - Delete All Collection
func (p *LeaveMongoDBDao) DeleteAll() (int64, error) {

	log.Println("accountMongoDao::DeleteAll - Begin ")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return 0, err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filter := bson.E{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId}

	res, err := collection.DeleteMany(ctx, filter, opts)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("accountMongoDao::DeleteAll - End deleted %v documents\n", res.DeletedCount)
	return res.DeletedCount, nil
}

func (p *LeaveMongoDBDao) appendListLookups(stages []bson.M) []bson.M {

	// Lookup Stage for PlatformAppUsers ========================================
	lookupStage := bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         platform_common.DbPlatformAppUsers,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_STAFF_ID,
			hr_common.MONGODB_STR_FOREIGNFIELD: platform_common.FLD_APP_USER_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_STAFF_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	// Lookup Stage for User ==========================================
	lookupStage = bson.M{
		hr_common.MONGODB_LOOKUP: bson.M{
			hr_common.MONGODB_STR_FROM:         hr_common.DbHrLeaveTypes,
			hr_common.MONGODB_STR_LOCALFIELD:   hr_common.FLD_LEAVETYPE_ID,
			hr_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_LEAVETYPE_ID,
			hr_common.MONGODB_STR_AS:           hr_common.FLD_LEAVE_INFO,
			hr_common.MONGODB_STR_PIPELINE: []bson.M{
				// Remove following fields from result-set
				{hr_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID:              0,
					platform_common.FLD_APP_USER_ID:       0,
					"auth_key":                            0,
					platform_common.FLD_APP_USER_PASSWORD: 0,
					db_common.FLD_IS_DELETED:              0,
					db_common.FLD_CREATED_AT:              0,
					db_common.FLD_UPDATED_AT:              0}},
			},
		},
	}
	// Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	return stages
}
