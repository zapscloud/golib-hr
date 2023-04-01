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
	GetDetails(leave_id string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(leave_id string, indata utils.Map) (utils.Map, error)
	Delete(leave_id string, delete_permanent bool) error

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
	child               LeaveService
	businessID          string
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

	// Assign the BusinessId
	p.businessID = businessId

	// Instantiate other services
	p.daoLeave = hr_repository.NewLeaveDao(p.GetClient(), p.businessID)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoPlatformBusiness.GetDetails(p.businessID)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business id", ErrorDetail: "Given business id is not exist"}
		return nil, err
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
func (p *leaveBaseService) GetDetails(leave_id string) (utils.Map, error) {
	log.Printf("AccountService::FindByCode::  Begin %v", leave_id)

	data, err := p.daoLeave.Get(leave_id)
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
	indata[hr_common.FLD_BUSINESS_ID] = p.businessID
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
func (p *leaveBaseService) Update(leave_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AccountService::Update - Begin")

	data, err := p.daoLeave.Get(leave_id)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, hr_common.FLD_LEAVE_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)

	data, err = p.daoLeave.Update(leave_id, indata)
	log.Println("AccountService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *leaveBaseService) Delete(leave_id string, delete_permanent bool) error {

	log.Println("AccountService::Delete - Begin", leave_id)

	daoLeave := p.daoLeave
	if delete_permanent {
		result, err := daoLeave.Delete(leave_id)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(leave_id, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("LeaveService::Delete - End")
	return nil
}
