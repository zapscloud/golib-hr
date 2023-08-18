package hr_common

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-platform/platform_common"
)

// HR Module tables
const (
	DbPrefix = db_common.DB_COLLECTION_PREFIX

	DbHrStaffs        = DbPrefix + "hr_staffs"
	DbHrStaffTypes    = DbPrefix + "hr_staff_types"
	DbHrDepartments   = DbPrefix + "hr_departments"
	DbHrDesignations  = DbPrefix + "hr_designations"
	DbHrPositions     = DbPrefix + "hr_positions"
	DbHrAttendances   = DbPrefix + "hr_attendances"
	DbHrLeaves        = DbPrefix + "hr_leaves"
	DbHrHolidays      = DbPrefix + "hr_holidays"
	DbHrShifts        = DbPrefix + "hr_shifts"
	DbHrWorkLocations = DbPrefix + "hr_work_locations"
	DbHrLeaveTypes    = DbPrefix + "hr_leave_types"
)

// HR Module table fields
const (
	// Common fields for all tables
	FLD_BUSINESS_ID = platform_common.FLD_BUSINESS_ID

	// Staff table fields
	FLD_STAFF_ID   = "staff_id"
	FLD_STAFF_DATA = "staff_data"
	FLD_STAF_INFO  = "staff_info"

	// StaffType table fields
	FLD_STAFFTYPE_ID          = "staff_type_id"
	FLD_STAFFTYPE_NAME        = "staff_type_name"
	FLD_STAFFTYPE_DESCRIPTION = "staff_type_description"

	// Leave Type table fields
	FLD_LEAVETYPE_ID   = "leave_type_id"
	FLD_LEAVETYPE_NAME = "leave_type_name"
	FLD_LEAVETYPE_DESC = "leave_type_desc"

	// Department table fields
	FLD_DEPARTMENT_ID   = "department_id"
	FLD_DEPARTMENT_NAME = "department_name"
	FLD_DEPARTMENT_DESC = "department_desc"

	// Holiday table fileds
	FLD_HOLIDAY_ID          = "holiday_id"
	FLD_HOLIDAY_NAME        = "holiday_name"
	FLD_HOLIDAY_DATE        = "holiday_date"
	FLD_HOLIDAY_DESCRIPTION = "holiday_description"

	// Designation table fields
	FLD_DESIGNATION_ID          = "designation_id"
	FLD_DESIGNATION_NAME        = "designation_name"
	FLD_DESIGNATION_DESCRIPTION = "designation_description"

	FLD_POSITION_ID   = "position_id"
	FLD_POSITION_NAME = "position_name"

	// Attendance Table
	FLD_ATTENDANCE_ID   = "attendance_id" // Auto generated
	FLD_ATTENDANCE_TYPE = "type"          // Possible values "IN", "OUT"
	FLD_DATETIME        = "date_time"
	FLD_LATITUDE        = "latitude"
	FLD_LONGITUDE       = "longitude"

	// Leave Table
	FLD_LEAVE_ID          = "leave_id"
	FLD_LEAVE_FROM        = "leave_from"
	FLD_LEAVE_TO          = "leave_to"
	FLD_LEAVE_DESCRIPTION = "leave_description"
	FLD_LEAVE_APPROVED    = "leave_approved"
	FLD_LEAVE_TYPE        = "leave_type"

	// Shift Table
	FLD_SHIFT_ID          = "shift_id"
	FLD_SHIFT_FROM        = "shift_from"
	FLD_SHIFT_TO          = "shift_to"
	FLD_SHIFT_DESCRIPTION = "shift_description"

	// Work Location Table
	FLD_WORKLOCATION_ID          = "work_location_id"
	FLD_WORKLOCATION_NAME        = "work_location_name"
	FLD_WORKLOCATION_DESCRIPTION = "work_location_description"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

}

func GetServiceModuleCode() string {
	return "HR"
}
