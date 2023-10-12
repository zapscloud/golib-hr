package hr_services

import (
	"log"
	"time"

	"github.com/zapscloud/golib-business/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

const (
// // The character encoding for the email.
// CharSet = "UTF-8"
)

// AttendanceService - Attendances Service structure
type AttendanceService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(attendance_id string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	// Create(indata utils.Map) (utils.Map, error)
	// CreateMany(indata utils.Map) (utils.Map, error)
	ClockIn(indata utils.Map) (utils.Map, error)
	ClockInMany(indata utils.Map) (utils.Map, error)
	ClockOut(attendance_id string, indata utils.Map) (utils.Map, error)
	ClockOutMany(indata utils.Map) (utils.Map, error)
	Update(attendance_id string, indata utils.Map) (utils.Map, error)
	Delete(attendance_id string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// LoyaltyCardService - Attendances Service structure
type attendanceBaseService struct {
	db_utils.DatabaseService
	daoAttendance       hr_repository.AttendanceDao
	daoPlatformBusiness platform_repository.BusinessDao
	daoStaff            hr_repository.StaffDao

	child      AttendanceService
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewAttendanceService(props utils.Map) (AttendanceService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := attendanceBaseService{}

	// Open Database Service
	err := p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return p.errorReturn(err)
	}

	// Verify whether the User id data passed, this is optional parameter
	staffId, _ := utils.GetMemberDataStr(props, hr_common.FLD_STAFF_ID)
	// if err != nil {
	// 	return p.errorReturn(err)
	// }

	// Assign the BusinessId & StaffId
	p.businessId = businessId
	p.staffId = staffId

	p.daoAttendance = hr_repository.NewAttendanceDao(p.GetClient(), p.businessId, p.staffId)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoStaff = hr_repository.NewStaffDao(p.GetClient(), p.businessId)

	_, err = p.daoPlatformBusiness.Get(p.businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business id", ErrorDetail: "Given business id is not exist"}
		return p.errorReturn(err)
	}

	// Verify the Staff Exist
	if len(staffId) > 0 {
		_, err = p.daoStaff.Get(staffId)
		if err != nil {
			err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid StaffId", ErrorDetail: "Given StaffId is not exist"}
			return p.errorReturn(err)
		}
	}

	p.child = &p

	return &p, nil
}

func (p *attendanceBaseService) EndService() {
	p.CloseDatabaseService()
}

// ************************
// List - List All records
//
// ************************
func (p *attendanceBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AttendanceService::FindAll - Begin")

	daoAttendance := p.daoAttendance
	response, err := daoAttendance.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AttendanceService::FindAll - End ")
	return response, nil
}

// *************************
// Get - Get Details
//
// *************************
func (p *attendanceBaseService) Get(appattendance_id string) (utils.Map, error) {
	log.Printf("AttendanceService::FindByCode::  Begin %v", appattendance_id)

	data, err := p.daoAttendance.Get(appattendance_id)
	log.Println("AttendanceService::FindByCode:: End ", err)
	return data, err
}

// ************************
// Find - Find Service
//
// ************************
func (p *attendanceBaseService) Find(filter string) (utils.Map, error) {
	log.Println("AttendanceService::FindByCode::  Begin ", filter)

	data, err := p.daoAttendance.Find(filter)
	log.Println("AttendanceService::FindByCode:: End ", data, err)
	return data, err
}

// ************************
// Create - Create Service
//
// ************************
// func (p *attendanceBaseService) Create(indata utils.Map) (utils.Map, error) {

// 	log.Println("AttendanceService::Create - Begin")

// 	// Create AttendanceId
// 	attendanceId := p.createAttendanceId(indata)

// 	if utils.IsEmpty(p.staffId) {
// 		err := &utils.AppError{
// 			ErrorCode:   "S30102",
// 			ErrorMsg:    "No StaffId",
// 			ErrorDetail: "No StaffId passed"}
// 		return indata, err
// 	}

// 	indata[hr_common.FLD_ATTENDANCE_ID] = attendanceId
// 	indata[hr_common.FLD_BUSINESS_ID] = p.businessId
// 	indata[hr_common.FLD_STAFF_ID] = p.staffId
// 	indata[hr_common.FLD_DATETIME] = time.Now().UTC() //.Format("2006-01-02 15:04:05")

// 	log.Println("Provided Attendance ID:", attendanceId)

// 	insertResult, err := p.daoAttendance.Create(indata)
// 	log.Println("AttendanceService::Create - End ", insertResult)

// 	return indata, err
// }

// ********************************
// CreateMany - CreateMany Service
//
// ********************************
// func (p *attendanceBaseService) CreateMany(indata utils.Map) (utils.Map, error) {

// 	var err error = nil

// 	log.Println("AttendanceService::CreateMany - Begin")

// 	// Create AttendanceId
// 	attendanceId := p.createAttendanceId(indata)

// 	// Check staffId received in indata
// 	staffId, _ := utils.GetMemberDataStr(indata, hr_common.FLD_STAFF_ID)
// 	if utils.IsEmpty(staffId) {
// 		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "No StaffId", ErrorDetail: "No StaffId passed"}
// 		return indata, err
// 	}
// 	indata[hr_common.FLD_ATTENDANCE_ID] = attendanceId
// 	indata[hr_common.FLD_BUSINESS_ID] = p.businessId

// 	// Convert Date_time string to Date Format
// 	if dataVal, dataOk := indata[hr_common.FLD_DATETIME]; dataOk {
// 		layout := hr_common.DATETIME_PARSE_FORMAT
// 		indata[hr_common.FLD_DATETIME], err = time.Parse(layout, dataVal.(string))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	insertResult, err := p.daoAttendance.Create(indata)

// 	log.Println("AttendanceService::Create - End ", insertResult)
// 	return indata, err
// }

// *************************
// ClockIn - Clock IN
//
// ************************
func (p *attendanceBaseService) ClockIn(indata utils.Map) (utils.Map, error) {
	var err error = nil

	log.Println("AttendanceService::ClockIn - Begin")

	// Create AttendanceId
	attendanceId := p.createAttendanceId(indata)

	// Add Current DateTime
	indata[hr_common.FLD_DATETIME] = time.Now().UTC()

	// Create ClockIn Data
	var clockIn utils.Map = utils.Map{}

	clockIn[hr_common.FLD_ATTENDANCE_ID] = attendanceId
	clockIn[hr_common.FLD_BUSINESS_ID] = p.businessId
	clockIn[hr_common.FLD_STAFF_ID] = p.staffId

	// Update Clock-In Interface back
	clockIn[hr_common.FLD_CLOCK_IN] = indata

	_, err = p.daoAttendance.Create(clockIn)

	log.Println("AttendanceService::ClockIn - End")
	return clockIn, err
}

// *************************************************
// ClockInMany - Clock In with StaffId and DateTime
//
// ************************************************
func (p *attendanceBaseService) ClockInMany(indata utils.Map) (utils.Map, error) {
	var err error = nil

	log.Println("AttendanceService::ClockInMany - Begin")

	// Create AttendanceId
	attendanceId := p.createAttendanceId(indata)

	// Check staffId received in indata
	staffId, _ := utils.GetMemberDataStr(indata, hr_common.FLD_STAFF_ID)
	_, err = p.daoStaff.Get(staffId)
	if err != nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Invalid StaffId", ErrorDetail: "No such StaffId found"}
		return indata, err
	}
	// Get TimeZone
	timeZone, err := utils.GetMemberDataStr(indata, business_common.FLD_BUSINESS_TIME_ZONE)
	if err != nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Need Timezone", ErrorDetail: "No TimeZone value passed"}
		return nil, err
	}

	// Convert Date_time string to Date Format
	dateTime, err := utils.GetMemberDataStr(indata, hr_common.FLD_DATETIME)
	if err == nil {
		// Load Location
		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			return nil, err
		}

		layout := hr_common.DATETIME_PARSE_FORMAT
		//indata[hr_common.FLD_DATETIME], err = time.Parse(layout, dateTime)
		indata[hr_common.FLD_DATETIME], err = time.ParseInLocation(layout, dateTime, loc)
		if err != nil {
			return nil, err
		}
	}
	// Remove StaffId from indata
	delete(indata, hr_common.FLD_STAFF_ID)
	// Remove Timezone from indata
	delete(indata, business_common.FLD_BUSINESS_TIME_ZONE)

	// Prepare ClockIn Data
	var clockIn utils.Map = utils.Map{}
	clockIn[hr_common.FLD_ATTENDANCE_ID] = attendanceId
	clockIn[hr_common.FLD_BUSINESS_ID] = p.businessId
	clockIn[hr_common.FLD_STAFF_ID] = staffId

	// Update Clock-In Interface back
	clockIn[hr_common.FLD_CLOCK_IN] = indata

	insertResult, err := p.daoAttendance.Create(clockIn)

	log.Println("AttendanceService::ClockInMany - End ", insertResult)
	return clockIn, err

}

// *************************
// ClockOut - Clock Out
//
// ************************
func (p *attendanceBaseService) ClockOut(attendance_id string, indata utils.Map) (utils.Map, error) {
	var err error = nil

	log.Println("AttendanceService::ClockOut - Begin")

	data, err := p.daoAttendance.Get(attendance_id)
	if err != nil {
		return indata, err
	}

	// Update DateTime
	indata[hr_common.FLD_DATETIME] = time.Now().UTC()

	// Update Clock-In Interface back
	data[hr_common.FLD_CLOCK_OUT] = indata

	_, err = p.daoAttendance.Update(attendance_id, data)

	log.Println("AttendanceService::ClockIn - End")
	return data, err
}

// *************************************************
// ClockInMany - Clock Out with AttendanceId and DateTime
//
// ************************************************
func (p *attendanceBaseService) ClockOutMany(indata utils.Map) (utils.Map, error) {
	var err error = nil

	log.Println("AttendanceService::ClockOutMany - Begin")
	// Check staffId received in indata
	attendanceId, _ := utils.GetMemberDataStr(indata, hr_common.FLD_ATTENDANCE_ID)
	data, err := p.daoAttendance.Get(attendanceId)
	if err != nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Invalid AttendanceId", ErrorDetail: "No such AttendanceId found"}
		return nil, err
	}
	// Get TimeZone
	timeZone, err := utils.GetMemberDataStr(indata, business_common.FLD_BUSINESS_TIME_ZONE)
	if err != nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Need Timezone", ErrorDetail: "No TimeZone value passed"}
		return nil, err
	}
	// Convert Date_time string to Date Format
	dateTime, err := utils.GetMemberDataStr(indata, hr_common.FLD_DATETIME)
	if err == nil {
		// Load Location
		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			return nil, err
		}
		layout := hr_common.DATETIME_PARSE_FORMAT
		//indata[hr_common.FLD_DATETIME], err = time.Parse(layout, dateTime)
		indata[hr_common.FLD_DATETIME], err = time.ParseInLocation(layout, dateTime, loc)
		if err != nil {
			return nil, err
		}
	}
	// Remove StaffId from indata
	delete(indata, hr_common.FLD_ATTENDANCE_ID)

	// Update Clock-In Interface back
	data[hr_common.FLD_CLOCK_OUT] = indata

	_, err = p.daoAttendance.Update(attendanceId, data)

	log.Println("AttendanceService::ClockIn - End")
	return data, err

}

// ************************
// Update - Update Service
//
// ************************
func (p *attendanceBaseService) Update(attendance_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AttendanceService::Update - Begin")

	data, err := p.daoAttendance.Get(attendance_id)
	if err != nil {
		return data, err
	}

	// Delete the Key fields
	delete(indata, hr_common.FLD_ATTENDANCE_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)
	delete(indata, hr_common.FLD_STAFF_ID)
	delete(indata, hr_common.FLD_DATETIME)

	// Convert clock_in->date_time string to Date Format
	if dataVal, dataOk := indata[hr_common.FLD_CLOCK_IN]; dataOk {
		_, err = p.convertStrToDateFormat(dataVal.(map[string]interface{}))
		if err != nil {
			log.Println("Failed to Parse clock_in->date_time", err)
			return nil, err
		}
	}

	// Convert clock_in->date_time string to Date Format
	if dataVal, dataOk := indata[hr_common.FLD_CLOCK_OUT]; dataOk {
		_, err = p.convertStrToDateFormat(dataVal.(map[string]interface{}))
		if err != nil {
			log.Println("Failed to Parse clock_out->date_time", err)
			return nil, err
		}
	}

	data, err = p.daoAttendance.Update(attendance_id, indata)
	log.Println("AttendanceService::Update - End ")
	return data, err
}

// ************************
// Delete - Delete Service
//
// ************************
func (p *attendanceBaseService) Delete(attendance_id string, delete_permanent bool) error {

	log.Println("AttendanceService::Delete - Begin", attendance_id, delete_permanent)

	daoAttendance := p.daoAttendance
	if delete_permanent {
		result, err := daoAttendance.Delete(attendance_id)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(attendance_id, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("AttendanceService::Delete - End")
	return nil
}

func (p *attendanceBaseService) errorReturn(err error) (AttendanceService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}

func (p *attendanceBaseService) convertStrToDateFormat(indata utils.Map) (utils.Map, error) {
	var err error = nil

	// Convert Date_time string to Date Format
	if dataVal, dataOk := indata[hr_common.FLD_DATETIME]; dataOk {
		layout := hr_common.DATETIME_PARSE_FORMAT
		indata[hr_common.FLD_DATETIME], err = time.Parse(layout, dataVal.(string))
		if err != nil {
			return nil, err
		}
	}

	return indata, err
}

func (p *attendanceBaseService) createAttendanceId(indata utils.Map) string {

	return utils.GenerateUniqueId("atten")
}
