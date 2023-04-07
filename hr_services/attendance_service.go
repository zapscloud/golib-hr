package hr_services

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	Create(indata utils.Map) (utils.Map, error)
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
	businessId, err := utils.IsMemberExist(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	// Verify whether the User id data passed, this is optional parameter
	staffId, _ := utils.IsMemberExist(props, hr_common.FLD_STAFF_ID)
	// if err != nil {
	// 	return nil, err
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
		return nil, err
	}

	// Verify the Staff Exist
	if len(staffId) > 0 {
		_, err = p.daoStaff.Get(staffId)
		if err != nil {
			err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid StaffId", ErrorDetail: "Given StaffId is not exist"}
			return nil, err
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
// GetDetails - Get Details
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
	fmt.Println("AttendanceService::FindByCode::  Begin ", filter)

	data, err := p.daoAttendance.Find(filter)
	log.Println("AttendanceService::FindByCode:: End ", data, err)
	return data, err
}

// ************************
// Create - Create Service
//
// ************************
func (p *attendanceBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("AttendanceService::Create - Begin")
	var attendanceId string

	dataval, dataok := indata[hr_common.FLD_ATTENDANCE_ID]
	if dataok {
		attendanceId = strings.ToLower(dataval.(string))
	} else {
		attendanceId = utils.GenerateUniqueId("atten")
		log.Println("Unique Attendance ID", attendanceId)
	}

	indata[hr_common.FLD_ATTENDANCE_ID] = attendanceId
	indata[hr_common.FLD_BUSINESS_ID] = p.businessId
	indata[hr_common.FLD_STAFF_ID] = p.staffId
	indata[hr_common.FLD_DATETIME] = time.Now()
	log.Println("Provided Attendance ID:", attendanceId)

	_, err := p.daoAttendance.Get(attendanceId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Attendance ID !", ErrorDetail: "Given Attendance ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoAttendance.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("AttendanceService::Create - End ", insertResult)
	return indata, err
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
