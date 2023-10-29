package hr_services

import (
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-platform/platform_services"
	"github.com/zapscloud/golib-utils/utils"
)

// HoursfactorService - Accounts Service structure
type HoursfactorService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(hoursfactorId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(hoursfactorId string, indata utils.Map) (utils.Map, error)
	Delete(hoursfactorId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// HoursfactoreBaseService - Accounts Service structure
type HoursfactorBaseService struct {
	db_utils.DatabaseService
	dbRegion            db_utils.DatabaseService
	daoHrsFactor        hr_repository.HoursfactorDao
	daoPlatformBusiness platform_repository.BusinessDao

	child      HoursfactorService
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewHoursfactorService(props utils.Map) (HoursfactorService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	log.Printf("HoursfactorService::Start ")

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	p := HoursfactorBaseService{}

	// Open Database Service
	err = p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Open RegionDB Service
	p.dbRegion, err = platform_services.OpenRegionDatabaseService(props)
	if err != nil {
		p.CloseDatabaseService()
		return nil, err
	}

	// Assign the BusinessId & StaffId
	p.businessId = businessId

	// Instantiate other services
	p.daoHrsFactor = hr_repository.NewHoursfactorDao(p.dbRegion.GetClient(), p.businessId)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoPlatformBusiness.Get(p.businessId)
	if err != nil {
		err := &utils.AppError{
			ErrorCode:   funcode + "01",
			ErrorMsg:    "Invalid business id",
			ErrorDetail: "Given business id is not exist"}
		return p.errorReturn(err)
	}

	p.child = &p

	return &p, nil
}

func (p *HoursfactorBaseService) EndService() {
	p.CloseDatabaseService()
	p.dbRegion.CloseDatabaseService()
}

// List - List All records
func (p *HoursfactorBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("HoursfactorService::FindAll - Begin")

	daoHrsFactor := p.daoHrsFactor
	response, err := daoHrsFactor.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("HoursfactorService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *HoursfactorBaseService) Get(hoursfactorId string) (utils.Map, error) {
	log.Printf("HoursfactorService::FindByCode::  Begin %v", hoursfactorId)

	data, err := p.daoHrsFactor.Get(hoursfactorId)
	log.Println("HoursfactorService::FindByCode:: End ", err)
	return data, err
}

func (p *HoursfactorBaseService) Find(filter string) (utils.Map, error) {
	log.Println("HoursfactorService::FindByCode::  Begin ", filter)

	data, err := p.daoHrsFactor.Find(filter)
	log.Println("HoursfactorService::FindByCode:: End ", data, err)
	return data, err
}

func (p *HoursfactorBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var hoursfactorId string

	dataval, dataok := indata[hr_common.FLD_HOURSFACTOR_ID]
	if dataok {
		hoursfactorId = strings.ToLower(dataval.(string))
	} else {
		hoursfactorId = utils.GenerateUniqueId("hfprof")
		log.Println("Unique Account ID", hoursfactorId)
	}
	indata[hr_common.FLD_HOURSFACTOR_ID] = hoursfactorId
	indata[hr_common.FLD_BUSINESS_ID] = p.businessId
	log.Println("Provided Account ID:", hoursfactorId)

	_, err := p.daoHrsFactor.Get(hoursfactorId)
	if err == nil {
		err := &utils.AppError{
			ErrorCode:   "S30102",
			ErrorMsg:    "Existing Hours Factor ID !",
			ErrorDetail: "Given Hours Factor ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoHrsFactor.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *HoursfactorBaseService) Update(hoursfactorId string, indata utils.Map) (utils.Map, error) {

	log.Println("HoursfactorService::Update - Begin")

	data, err := p.daoHrsFactor.Get(hoursfactorId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, hr_common.FLD_HOURSFACTOR_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)

	data, err = p.daoHrsFactor.Update(hoursfactorId, indata)
	log.Println("HoursfactorService::Update - End ", err)
	return data, err
}

// Delete - Delete Service
func (p *HoursfactorBaseService) Delete(hoursfactorId string, delete_permanent bool) error {

	log.Println("HoursfactorService::Delete - Begin", hoursfactorId)

	daoHrsFactor := p.daoHrsFactor
	if delete_permanent {
		result, err := daoHrsFactor.Delete(hoursfactorId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(hoursfactorId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("HoursfactorService::Delete - End")
	return nil
}

func (p *HoursfactorBaseService) errorReturn(err error) (HoursfactorService, error) {
	// Close the Database Connection
	p.EndService()
	return nil, err
}
