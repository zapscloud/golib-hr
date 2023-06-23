package hr_services

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// StaffService - Accounts Service structure
type StaffService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(staff_id string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(staff_id string, indata utils.Map) (utils.Map, error)
	Delete(staff_id string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// staffBaseService - Accounts Service structure
type staffBaseService struct {
	db_utils.DatabaseService
	daoStaff            hr_repository.StaffDao
	daoPlatformBusiness platform_repository.BusinessDao
	child               StaffService
	businessID          string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewStaffService(props utils.Map) (StaffService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := staffBaseService{}

	// Open Database Service
	err := p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	// Assign the BusinessId
	p.businessID = businessId

	// Instantiate other services
	p.daoStaff = hr_repository.NewStaffDao(p.GetClient(), p.businessID)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoPlatformBusiness.Get(p.businessID)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business id", ErrorDetail: "Given business id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, nil
}

func (p *staffBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *staffBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AccountService::FindAll - Begin")

	daoStaff := p.daoStaff
	response, err := daoStaff.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AccountService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *staffBaseService) Get(staff_id string) (utils.Map, error) {
	log.Printf("AccountService::FindByCode::  Begin %v", staff_id)

	data, err := p.daoStaff.Get(staff_id)
	log.Println("AccountService::FindByCode:: End ", err)
	return data, err
}

func (p *staffBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AccountService::FindByCode::  Begin ", filter)

	data, err := p.daoStaff.Find(filter)
	log.Println("AccountService::FindByCode:: End ", data, err)
	return data, err
}

func (p *staffBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	dataval, dataok := indata[hr_common.FLD_STAFF_ID]
	if !dataok {
		uid := utils.GenerateUniqueId("stf")
		log.Println("Unique Account ID", uid)
		indata[hr_common.FLD_STAFF_ID] = uid
		dataval = indata[hr_common.FLD_STAFF_ID]
	}
	indata[hr_common.FLD_BUSINESS_ID] = p.businessID
	log.Println("Provided Account ID:", dataval)

	_, err := p.daoStaff.Get(dataval.(string))
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Account ID !", ErrorDetail: "Given Account ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoStaff.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *staffBaseService) Update(staff_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AccountService::Update - Begin")

	data, err := p.daoStaff.Get(staff_id)
	if err != nil {
		return data, err
	}

	data, err = p.daoStaff.Update(staff_id, indata)
	log.Println("AccountService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *staffBaseService) Delete(staff_id string, delete_permanent bool) error {

	log.Println("AccountService::Delete - Begin", staff_id)

	daoStaff := p.daoStaff
	if delete_permanent {
		result, err := daoStaff.Delete(staff_id)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(staff_id, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("StaffService::Delete - End")
	return nil
}
