package hr_services

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_repository"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-platform/platform_services"
	"github.com/zapscloud/golib-utils/utils"
)

// DashboardService - Dashboard Service structure
type DashboardService interface {
	GetDashboardData() (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type dashboardBaseService struct {
	db_utils.DatabaseService
	dbRegion     db_utils.DatabaseService
	daoDashboard hr_repository.DashboardDao
	daoBusiness  platform_repository.BusinessDao
	child        DashboardService
	businessID   string
	staffID      string // Changed "staffId" to "staffID" for consistency
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewDashboardService(props utils.Map) (DashboardService, error) {
	funcode := hr_common.GetServiceModuleCode() + "M" + "01"

	log.Printf("DashboardService :: Start")

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, hr_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	p := dashboardBaseService{}

	// Open Database Service
	err = p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Open RegionDB Service
	err = p.openRegionDatabaseService(props)
	if err != nil {
		p.CloseDatabaseService()
		return nil, err
	}

	// Verify whether the User id data passed, this is optional parameter
	staffID, _ := utils.GetMemberDataStr(props, hr_common.FLD_STAFF_ID)

	// Assign the BusinessID
	p.businessID = businessId
	p.staffID = staffID

	// Instantiate other services
	p.daoDashboard = hr_repository.NewDashboardDao(p.dbRegion.GetClient(), p.businessID, p.staffID)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{
			ErrorCode:   funcode + "01",
			ErrorMsg:    "Invalid business_id",
			ErrorDetail: "Given app_business_id does not exist"}
		return p.errorReturn(err)
	}

	// // Verify the Staff Exist
	// if len(staffID) > 0 {
	// 	_, err = p.daoStaff.Get(staffID)
	// 	if err != nil {
	// 		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid StaffId", ErrorDetail: "Given StaffId is not exist"}
	// 		return p.errorReturn(err)
	// 	}
	// }

	p.child = &p

	return &p, err
}

func (p *dashboardBaseService) EndService() {
	log.Printf("EndDashboardMongoService ")
	p.CloseDatabaseService()
	p.dbRegion.CloseDatabaseService()
}

func (p *dashboardBaseService) openRegionDatabaseService(props utils.Map) error {

	// Get Region and Tenant Database Information
	propsRegion, err := platform_services.GetRegionAndTenantDBInfo(props)
	if err != nil {
		log.Println("GetRegionAndTenantDBInfo() ERROR", err)
		return err
	}

	err = p.dbRegion.OpenDatabaseService(propsRegion)
	if err != nil {
		return err
	}

	return nil
}

// GetDashboardData retrieves dashboard data
func (p *dashboardBaseService) GetDashboardData() (utils.Map, error) {
	log.Println("DashboardService::GetDashboardData - Begin")

	daoDashboard := p.daoDashboard
	response, err := daoDashboard.GetDashboardData()
	if err != nil {
		return nil, err
	}

	log.Println("DashboardService::GetDashboardData - End")
	return response, nil
}

// errorReturn handles error and closes the database connection
func (p *dashboardBaseService) errorReturn(err error) (DashboardService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
