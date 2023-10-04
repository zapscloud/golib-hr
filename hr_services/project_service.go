package hr_services

import (
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// ProjectService - Accounts Service structure
type ProjectService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(projectId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(projectId string, indata utils.Map) (utils.Map, error)
	Delete(projectId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// ProjectBaseService - Accounts Service structure
type ProjectBaseService struct {
	db_utils.DatabaseService
	daoProject          hr_repository.ProjectDao
	daoPlatformBusiness platform_repository.BusinessDao
	child               ProjectService
	businessID          string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewProjectService(props utils.Map) (ProjectService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := ProjectBaseService{}

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

	// Assign the BusinessId
	p.businessID = businessId

	// Instantiate other services
	p.daoProject = hr_repository.NewProjectDao(p.GetClient(), p.businessID)
	p.daoPlatformBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoPlatformBusiness.Get(p.businessID)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business id", ErrorDetail: "Given business id is not exist"}
		return p.errorReturn(err)
	}

	p.child = &p

	return &p, nil
}

func (p *ProjectBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *ProjectBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AccountService::FindAll - Begin")

	daoProject := p.daoProject
	response, err := daoProject.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AccountService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *ProjectBaseService) Get(projectId string) (utils.Map, error) {
	log.Printf("AccountService::FindByCode::  Begin %v", projectId)

	data, err := p.daoProject.Get(projectId)
	log.Println("AccountService::FindByCode:: End ", err)
	return data, err
}

func (p *ProjectBaseService) Find(filter string) (utils.Map, error) {
	log.Println("AccountService::FindByCode::  Begin ", filter)

	data, err := p.daoProject.Find(filter)
	log.Println("AccountService::FindByCode:: End ", data, err)
	return data, err
}

func (p *ProjectBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var ProjectId string

	dataval, dataok := indata[hr_common.FLD_PROJECT_ID]
	if dataok {
		ProjectId = strings.ToLower(dataval.(string))
	} else {
		ProjectId = utils.GenerateUniqueId("projt")
		log.Println("Unique Account ID", ProjectId)
	}
	indata[hr_common.FLD_PROJECT_ID] = ProjectId
	indata[hr_common.FLD_BUSINESS_ID] = p.businessID
	log.Println("Provided Account ID:", ProjectId)

	_, err := p.daoProject.Get(ProjectId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Account ID !", ErrorDetail: "Given Account ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoProject.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *ProjectBaseService) Update(projectId string, indata utils.Map) (utils.Map, error) {

	log.Println("AccountService::Update - Begin")

	data, err := p.daoProject.Get(projectId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, hr_common.FLD_PROJECT_ID)
	delete(indata, hr_common.FLD_BUSINESS_ID)

	data, err = p.daoProject.Update(projectId, indata)
	log.Println("AccountService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *ProjectBaseService) Delete(projectId string, delete_permanent bool) error {

	log.Println("AccountService::Delete - Begin", projectId)

	daoProject := p.daoProject
	if delete_permanent {
		result, err := daoProject.Delete(projectId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(projectId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("ProjectService::Delete - End")
	return nil
}

func (p *ProjectBaseService) errorReturn(err error) (ProjectService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
