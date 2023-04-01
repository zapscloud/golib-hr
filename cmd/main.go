package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr/hr_common"
	"github.com/zapscloud/golib-hr/hr_services"
	"github.com/zapscloud/golib-utils/utils"
)

func GetDBCreds() utils.Map {
	dbtype := db_common.DATABASE_TYPE_MONGODB
	dbuser := os.Getenv("MONGO_DB_USER")
	dbsecret := os.Getenv("MONGO_DB_SECRET")
	dbserver := os.Getenv("MONGO_DB_SERVER")
	dbname := os.Getenv("MONGO_DB_NAME")

	dbCreds := utils.Map{
		db_common.DB_TYPE:   dbtype,
		db_common.DB_SERVER: dbserver,
		db_common.DB_NAME:   dbname,
		db_common.DB_USER:   dbuser,
		db_common.DB_SECRET: dbsecret}

	return dbCreds
}

func MdbMain(businessid string) hr_services.AttendanceService {

	dbCreds := GetDBCreds()

	fmt.Println("DB Credentials: ", dbCreds)

	if dbCreds[db_common.DB_SERVER].(string) == "" {
		fmt.Println("Environment variable MONGO_DB_SERVER should be defined")
		return nil
	} else if dbCreds[db_common.DB_NAME].(string) == "" {
		fmt.Println("Environment variable MONGO_DB_NAME should be defined")
		return nil
	}

	dbCreds[hr_common.FLD_BUSINESS_ID] = businessid

	rolesrv, err := hr_services.NewAttendanceService(dbCreds)
	fmt.Println("User Mongo Service Error ", err)
	return rolesrv
}

func main() {

	businessid := "business003"
	rolesrv := MdbMain(businessid)
	// usersrv, bizsrv, rolesrv := ZapsMain(businessid)

	// EmptyBusiness(bizsrv)
	// DeleteBusiness(bizsrv)
	// CreateBusiness(bizsrv)
	// GetBusiness(bizsrv)

	if rolesrv != nil {
		EmptyBusinessAttendance(rolesrv)
		// DeleteAttendance(rolesrv)
		// CreateAttendance(rolesrv)
		// UpdateAttendance(rolesrv)
		ListAttendances(rolesrv)
		// GetAttendance(rolesrv)
		// FindAttendance(rolesrv)
	}

}

func EmptyBusinessAttendance(srv hr_services.AttendanceService) {
	fmt.Println("Attendance Service ")
}

func CreateAttendance(srv hr_services.AttendanceService) {

	indata := utils.Map{
		"role_id":    "role003",
		"role_name":  "Demo Attendance 003",
		"role_scope": "admin",
	}

	res, err := srv.Create(indata)
	fmt.Println("Create Attendance", err)
	pretty.Println(res)

}

func GetAttendance(srv hr_services.AttendanceService) {
	res, err := srv.Get("role001")
	fmt.Println("Get Attendance", err)
	pretty.Println(res)

}

func FindAttendance(srv hr_services.AttendanceService) {

	filter := fmt.Sprintf(`{"%s":"%s"}`, "role_scope", "admin")
	res, err := srv.Find(filter)
	fmt.Println("Get Attendance", err)
	pretty.Println(res)

}

func UpdateAttendance(srv hr_services.AttendanceService) {

	indata := utils.Map{
		"role_id":   "role001",
		"role_name": "Demo Attendance 001 Updated",
		"is_active": true,
	}

	res, err := srv.Update("role001", indata)
	fmt.Println("Update Attendance", err)
	pretty.Println(res)

}

func DeleteAttendance(srv hr_services.AttendanceService) {

	srv.BeginTransaction()
	err := srv.Delete("role001", false)
	fmt.Println("DeleteAttendance success ", err)
	fmt.Println("DeleteAttendance Value ")

	if err != nil {
		srv.RollbackTransaction()
	} else {
		srv.CommitTransaction()
	}
}

func ListAttendances(srv hr_services.AttendanceService) {

	filter := "" //fmt.Sprintf(`{"%s":"%s"}`, "role_scope", "admin")

	sort := `{ "role_scope":1, "role_id":1}`

	res, err := srv.List(filter, sort, 0, 0)
	fmt.Println("List User success ", err)
	fmt.Println("List User summary ", res)
	pretty.Print(res)
}
