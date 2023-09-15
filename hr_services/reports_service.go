package hr_services

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// ReportsService - Reports Service structure
type ReportsService interface {
	GetAttendanceSummary(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type reportsBaseService struct {
	db_utils.DatabaseService
	daoReports  hr_repository.ReportsDao
	daoBusiness platform_repository.BusinessDao

	child      ReportsService
	businessID string
	staffID    string // Changed "staffId" to "staffID" for consistency
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewReportsService(props utils.Map) (ReportsService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	p := reportsBaseService{} // Initialize p as a pointer to the struct

	// Open Database Service
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ReportsService ")

	// Verify whether the business id data passed
	businessID, err := utils.GetMemberDataStr(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return p.errorReturn(err)
	}

	// Verify whether the User id data passed, this is optional parameter
	staffID, _ := utils.GetMemberDataStr(props, hr_common.FLD_STAFF_ID)

	// Assign the BusinessID
	p.businessID = businessID
	p.staffID = staffID

	// Instantiate other services
	p.daoReports = hr_repository.NewReportsDao(p.GetClient(), p.businessID, p.staffID)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoBusiness.Get(businessID)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given app_business_id does not exist"}
		return p.errorReturn(err)
	}

	p.child = &p // Assign the pointer to itself

	return &p, err
}

func (p *reportsBaseService) EndService() {
	log.Printf("EndReportsMongoService ")
	p.CloseDatabaseService()
}

// GetAttendanceSummary retrieves reports data
func (p *reportsBaseService) GetAttendanceSummary(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	log.Println("ReportsService::GetReportsData - Begin")

	daoReports := p.daoReports
	response, err := daoReports.GetAttendanceSummary(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("ReportsService::GetAttendanceSummary - End")
	return response, nil
}

// errorReturn handles error and closes the database connection
func (p *reportsBaseService) errorReturn(err error) (ReportsService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
