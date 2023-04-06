package hr_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// LeaveService - Accounts Service structure
type LeaveService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	GetDetails(leaveId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(leaveId string, indata utils.Map) (utils.Map, error)
	Delete(leaveId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// leaveBaseService - Accounts Service structure
type leaveBaseService struct {
	db_utils.DatabaseService
	daoLeave            hr_repository.LeaveDao
	daoPlatformBusiness platform_repository.BusinessDao
	daoStaff            hr_repository.StaffDao

	child      LeaveService
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewLeaveService(props utils.Map) (LeaveService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := leaveBaseService{}

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

	// Instantiate other services
	p.daoLeave = hr_repository.NewLeaveDao(p.GetClient(), p.businessId, p.staffId)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoStaff = hr_repository.NewStaffDao(p.GetClient(), p.businessId)

	_, err = p.daoPlatformBusiness.GetDetails(p.businessId)
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

func (p *leaveBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *leaveBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AccountService::FindAll - Begin")

	daoLeave := p.daoLeave
	response, err := daoLeave.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AccountService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *leaveBaseService) GetDetails(leaveId string) (utils.Map, error) {
	log.Printf("AccountService::FindByCode::  Begin %v", leaveId)

	data, err := p.daoLeave.Get(leaveId)
	log.Println("AccountService::FindByCode:: End ", err)
	return data, err
}

func (p *leaveBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AccountService::FindByCode::  Begin ", filter)

	data, err := p.daoLeave.Find(filter)
	log.Println("AccountService::FindByCode:: End ", data, err)
	return data, err
}

func (p *leaveBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var leaveId string

	dataval, dataok := indata[hr_common.FLD_LEAVE_ID]
	if dataok {
		leaveId = strings.ToLower(dataval.(string))
	} else {
		leaveId = utils.GenerateUniqueId("leav")
		log.Println("Unique Account ID", leaveId)
	}
	indata[hr_common.FLD_LEAVE_ID] = leaveId
	indata[hr_common.FLD_BUSINESS_ID] = p.businessId
	indata[hr_common.FLD_STAFF_ID] = p.staffId
	log.Println("Provided Account ID:", leaveId)

	_, err := p.daoLeave.Get(leaveId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Account ID !", ErrorDetail: "Given Account ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoLeave.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *leaveBaseService) Update(leaveId string, indata utils.Map) (utils.Map, error) {

	log.Println("AccountService::Update - Begin")

	data, err := p.daoLeave.Get(leaveId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, hr_common.FLD_LEAVE_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)
	delete(indata, hr_common.FLD_STAFF_ID)

	data, err = p.daoLeave.Update(leaveId, indata)
	log.Println("AccountService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *leaveBaseService) Delete(leaveId string, delete_permanent bool) error {

	log.Println("AccountService::Delete - Begin", leaveId)

	daoLeave := p.daoLeave
	if delete_permanent {
		result, err := daoLeave.Delete(leaveId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(leaveId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("LeaveService::Delete - End")
	return nil
}
