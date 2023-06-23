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

// ShiftService - Accounts Service structure
type ShiftService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(shiftId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(shiftId string, indata utils.Map) (utils.Map, error)
	Delete(shiftId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// shiftBaseService - Accounts Service structure
type shiftBaseService struct {
	db_utils.DatabaseService
	daoShift            hr_repository.ShiftDao
	daoPlatformBusiness platform_repository.BusinessDao

	child      ShiftService
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewShiftService(props utils.Map) (ShiftService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := shiftBaseService{}

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

	// Assign the BusinessId & StaffId
	p.businessId = businessId

	// Instantiate other services
	p.daoShift = hr_repository.NewShiftDao(p.GetClient(), p.businessId)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoPlatformBusiness.Get(p.businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business id", ErrorDetail: "Given business id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, nil
}

func (p *shiftBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *shiftBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AccountService::FindAll - Begin")

	daoShift := p.daoShift
	response, err := daoShift.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AccountService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *shiftBaseService) Get(shiftId string) (utils.Map, error) {
	log.Printf("AccountService::FindByCode::  Begin %v", shiftId)

	data, err := p.daoShift.Get(shiftId)
	log.Println("AccountService::FindByCode:: End ", err)
	return data, err
}

func (p *shiftBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AccountService::FindByCode::  Begin ", filter)

	data, err := p.daoShift.Find(filter)
	log.Println("AccountService::FindByCode:: End ", data, err)
	return data, err
}

func (p *shiftBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var shiftId string

	dataval, dataok := indata[hr_common.FLD_SHIFT_ID]
	if dataok {
		shiftId = strings.ToLower(dataval.(string))
	} else {
		shiftId = utils.GenerateUniqueId("shift")
		log.Println("Unique Account ID", shiftId)
	}
	indata[hr_common.FLD_SHIFT_ID] = shiftId
	indata[hr_common.FLD_BUSINESS_ID] = p.businessId
	log.Println("Provided Account ID:", shiftId)

	_, err := p.daoShift.Get(shiftId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Account ID !", ErrorDetail: "Given Account ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoShift.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *shiftBaseService) Update(shiftId string, indata utils.Map) (utils.Map, error) {

	log.Println("AccountService::Update - Begin")

	data, err := p.daoShift.Get(shiftId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, hr_common.FLD_SHIFT_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)

	data, err = p.daoShift.Update(shiftId, indata)
	log.Println("AccountService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *shiftBaseService) Delete(shiftId string, delete_permanent bool) error {

	log.Println("AccountService::Delete - Begin", shiftId)

	daoShift := p.daoShift
	if delete_permanent {
		result, err := daoShift.Delete(shiftId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(shiftId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("ShiftService::Delete - End")
	return nil
}
