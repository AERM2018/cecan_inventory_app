package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitDepartmentsRoutes(router router.Party, dbPsql *gorm.DB) {
	departments := router.Party("/departments")
	departmentsDataSource := datasources.DepartmentDataSource{DbPsql: dbPsql}
	departmentsController := controllers.DepartmentsController{DepartmentsDataSource: departmentsDataSource}
	departmentsController.New()
	departments.Get("/", departmentsController.GetDepartments)
	departments.Post("/", departmentsController.CreateDepartment)
	departments.Put("/{id:string}", departmentsController.UpdateDepartment)
	departments.Put("/{id:string}/users/{userId:string}", departmentsController.AssingResponsibleToDepartment)
	departments.Put("/{id:string}/reactivate", departmentsController.ReactivateDepartment)
	departments.Delete("/{id:string}", departmentsController.DeleteDepartment)
}
