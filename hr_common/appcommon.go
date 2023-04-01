package hr_common

import (
	"log"
)

// HR Module tables
const (
	DbHrStaffs       = "hr_staffs"
	DbHrDepartments  = "hr_departments"
	DbHrDesignations = "hr_designations"
	DbHrPositions    = "hr_positions"
	DbHrAttendances  = "hr_attendances"
	DbHrLeaves       = "hr_leaves"
	DbHrHolidays     = "hr_holidays"
)

// HR Module table fields
const (
	// Common fields for all tables
	FLD_BUSINESS_ID = "business_id"

	// Staff table fields
	FLD_STAFF_ID = "staff_id"

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
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

}

func GetServiceModuleCode() string {
	return "HR"
}
