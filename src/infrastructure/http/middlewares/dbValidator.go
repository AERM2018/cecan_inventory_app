package middlewares

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type DbValidator struct {
	MedicineDataSrc               datasources.MedicinesDataSource
	PharmacyDataSrc               datasources.PharmacyStocksDataSource
	RolesDataSource               datasources.RolesDataSource
	UserDataSource                datasources.UserDataSource
	PrescriptionDataSource        datasources.PrescriptionsDataSource
	StorehouseUtilityDataSource   datasources.StorehouseUtilitiesDataSource
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
	StorehouseRequestsDataSource  datasources.StorehouseRequestsDataSource
	FixedAssetsDataSource         datasources.FixedAssetsDataSource
	FixedAssetsRequestsDataSource datasources.FixedAssetsRequetsDataSource
}

func (dbVal DbValidator) IsRoleId(ctx iris.Context) {
	var (
		httpRes models.Responser
		user    models.User
	)
	bodyreader.ReadBodyAsJson(ctx, &user, false)
	_, err := dbVal.RolesDataSource.GetRoleById(user.RoleId)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El rol con id: %v no existe.", user.RoleId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	}
	ctx.Next()
}

func (dbVal DbValidator) CanUserDoAction(roleNamesAllowed ...string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		var httpRes models.Responser
		// This is just for testing, include Admin role in all request
		roleNamesAllowed = append(roleNamesAllowed, "Admin")
		roleNamesAllowed = append(roleNamesAllowed, "Director")
		roleNamesAllowed = append(roleNamesAllowed, "Subdirector")
		roleName := ctx.Values().GetString("roleName")
		if !common.FindElementInSlice(roleName, roleNamesAllowed) {
			httpRes = models.Responser{
				StatusCode: iris.StatusForbidden,
				Message:    "Acción denegada, no cuenta con los permisos necesarios.",
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
		ctx.Next()
	}
}
func (dbVal DbValidator) IsEmail(ctx iris.Context) {
	var (
		httpRes          models.Responser
		user             models.User
		username         string
		credentialsReset models.AccessCredentialsRestart
		errNotFound      error
	)
	isPasswordReset := strings.Contains(ctx.FullRequestURI(), "password_reset_code")
	if isPasswordReset {
		bodyreader.ReadBodyAsJson(ctx, &credentialsReset, false)
		username = credentialsReset.Email

	} else {
		bodyreader.ReadBodyAsJson(ctx, &user, false)
		username = user.Email
	}
	_, errNotFound = dbVal.UserDataSource.GetUserByEmailOrId(username)
	if errNotFound == nil && !isPasswordReset {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El email %v ya está siendo usado por otro usuario.", user.Email),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	} else if errNotFound != nil && isPasswordReset {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El email %v no está asociado a ningún usuario.", credentialsReset.Email),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineInCatalogByKey(ctx iris.Context) {
	var (
		httpRes               models.Responser
		medicine              models.Medicine
		pharmacyStockToUpdate models.PharmacyStockToUpdate
		medicineKey           string
	)
	routePath := ctx.Path()
	if strings.Contains(routePath, "pharmacy_inventory") {
		if ctx.Request().Method == "PUT" {
			bodyreader.ReadBodyAsJson(ctx, &pharmacyStockToUpdate, false)
			medicineKey = pharmacyStockToUpdate.MedicineKey
		} else if ctx.Request().Method == "POST" {
			medicineKey = ctx.Params().GetString("key")
		}
	} else {
		medicineKey = ctx.Params().GetString("key")
		if medicineKey == "" { // Get medicine key from body when it's not found in the url
			bodyreader.ReadBodyAsJson(ctx, &medicine, false)
			medicineKey = medicine.Key
		}
	}
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Unscoped().Where("key = ?", medicineKey).First(&medicine).RowsAffected
	if isMedicine == 0 {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El medicamento con clave: %v no existe.", medicineKey),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineWithName(ctx iris.Context) {
	var (
		httpRes     models.Responser
		medicine    models.Medicine
		medicineRes models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	bodyreader.ReadBodyAsJson(ctx, &medicine, false)
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Unscoped().Where("key != ? AND name = ?", medicineKey, medicine.Name).Find(&medicineRes).RowsAffected
	if isMedicine != 0 {
		var (
			resMessage string
		)
		if ctx.Request().Method == "PUT" {
			resMessage = fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con el nombre: %v.", medicine.Name)
		} else {
			resMessage = fmt.Sprintf("Ya existe un medicamento con el nombre: %v.", medicine.Name)
		}
		if medicineRes.DeletedAt.Valid {
			resMessage += " (DESABILITADO)"
		}

		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    resMessage,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineWithKey(ctx iris.Context) {
	var (
		httpRes  models.Responser
		medicine models.Medicine
	)
	bodyreader.ReadBodyAsJson(ctx, &medicine, false)
	medicineKey := ctx.Params().GetString("key")
	if medicine.Key != medicineKey {
		var resMessage string
		isMedicineWithKey := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicine.Key).Find(&models.Medicine{}).RowsAffected
		if isMedicineWithKey == 1 {
			if ctx.Request().Method == "PUT" {
				resMessage = fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con la clave: %v.", medicine.Key)
			} else {
				resMessage = fmt.Sprintf("El medicamento con clave: %v ya existe.", medicine.Key)
			}
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    resMessage,
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineDeleted(ctx iris.Context) {
	var (
		httpRes models.Responser
		isError bool
		me      models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	reqPath := ctx.Path()
	IsMedicineDeleted := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicineKey).First(&me).RowsAffected
	var message string
	if strings.Contains(reqPath, "pharmacy_inventory") && ctx.Request().Method == "POST" && IsMedicineDeleted == 0 {
		isError = true
		message = fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", medicineKey)
	} else if strings.Contains(reqPath, "reactivate") && ctx.Request().Method == "PUT" && IsMedicineDeleted != 0 {
		isError = true
		message = fmt.Sprintf("El medicamento con clave: %v no se reactivó debido a que no ha sido eliminado antes.", medicineKey)
	}
	if isError {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    message,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPharmacyStockUsed(ctx iris.Context) {
	var (
		httpRes models.Responser
		action  string
	)
	if ctx.Method() == "PUT" {
		action = "actualizar"
	} else {
		action = "eliminar"
	}
	pharmacyStockId := ctx.Params().GetString("id")
	pharmacyStockUuid, _ := uuid.Parse(pharmacyStockId)
	isStockUsed, _ := dbVal.PharmacyDataSrc.IsStockUsed(pharmacyStockUuid)
	if isStockUsed {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("No se puede %v un stock de farmacia cuando ha sido utilizado.", action),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPharmacyStockById(ctx iris.Context) {
	var httpRes models.Responser
	pharmacyStockId := ctx.Params().GetString("id")
	pharmacyStockUuid, _ := uuid.Parse(pharmacyStockId)
	_, err := dbVal.PharmacyDataSrc.GetPharmacyStockById(pharmacyStockUuid)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El stock de farmacia con id: %v no existe.", pharmacyStockId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPharmacyStockWithLotNumber(ctx iris.Context) {
	var (
		pharmacyStockIdForUpdate string
		pharmacyStock            models.PharmacyStock
		httpRes                  models.Responser
		errorMessage             string
	)
	bodyreader.ReadBodyAsJson(ctx, &pharmacyStock, false)
	if ctx.Method() == "PUT" {
		pharmacyStockIdForUpdate = ctx.Params().GetStringDefault("id", "")
		errorMessage = fmt.Sprintf("No se pudó actualizar el stock en farmacia debido a que ya existe uno con el número de lote %v", pharmacyStock.LotNumber)
	} else {
		errorMessage = fmt.Sprintf("Ya existe un stock en farmacia con el número de lote %v", pharmacyStock.LotNumber)
	}
	isPharmacyStock := dbVal.PharmacyDataSrc.IsPharmacyStockWithLotNumber(pharmacyStock.LotNumber, pharmacyStockIdForUpdate)
	if isPharmacyStock {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    errorMessage,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPrescriptionDeterminedStatus(status string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		var httpRes models.Responser
		prescriptionId := ctx.Params().GetString("id")
		if !dbVal.PrescriptionDataSource.IsPrescriptionDeterminedStatus(prescriptionId, status) {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    fmt.Sprintf("No se pudó completar la acción, la receta no tiene un estado: %v", status),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
		ctx.Next()
	}
}

func (dbVal DbValidator) IsPrescriptionById(ctx iris.Context) {
	var httpRes models.Responser
	prescriptionId := ctx.Params().GetString("id")
	idUuid, _ := uuid.Parse(prescriptionId)
	_, err := dbVal.PrescriptionDataSource.GetPrescriptionById(idUuid)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("La receta con id: %v no existe.", prescriptionId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsSamePrescriptionCreator(ctx iris.Context) {
	var httpRes models.Responser
	prescriptionId := ctx.Params().GetString("id")
	creatorUserId := ctx.Values().GetString("userId")
	isSameCreator := dbVal.PrescriptionDataSource.IsSamePrescriptionCreator(prescriptionId, creatorUserId)
	if !isSameCreator {
		httpRes = models.Responser{
			StatusCode: iris.StatusForbidden,
			Message:    "Solo el creador de la receta está permitido a actualizarla/borrarla.",
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseUtilityWithKey(ctx iris.Context) {
	// If it's found, avoid repeatition
	var (
		httpRes models.Responser
		utility models.StorehouseUtility
	)
	bodyreader.ReadBodyAsJson(ctx, &utility, false)
	utilityKey := ctx.Params().GetString("key")
	if utility.Key != utilityKey {
		var resMessage string
		storehouseUtility, err := dbVal.StorehouseUtilityDataSource.GetStorehouseUtilityByKey(utility.Key)
		if err == nil {
			if ctx.Request().Method == "PUT" {
				resMessage = fmt.Sprintf("No se actualizó el elemento de almacen debido a que ya existe un elemento con la clave: %v.", utility.Key)
			} else {
				resMessage = fmt.Sprintf("El elemento de almacen con clave: %v ya existe.", utility.Key)
			}
			if storehouseUtility.DeletedAt.Valid {
				resMessage += " (DESABILITADO)"
			}
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    resMessage,
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) FindStorehouseUtilityByKey(ctx iris.Context) {
	// Look for it to know if the key specified exist in order to alterate the utility or a stock
	var (
		httpRes      models.Responser
		utility      models.StorehouseUtility
		utilityStock models.StorehouseStock
		utilityKey   string
	)
	//Cases when the key is in the params
	if (!strings.Contains(ctx.FullRequestURI(), "storehouse_inventory") && (ctx.Method() == "PUT" || ctx.Method() == "DELETE")) ||
		(strings.Contains(ctx.FullRequestURI(), "storehouse_inventory") && ctx.Method() == "POST") {
		utilityKey = ctx.Params().GetString("key")
		// The key is the struct of a storehouse stock
	} else if strings.Contains(ctx.FullRequestURI(), "storehouse_inventory") && ctx.Method() == "PUT" {
		bodyreader.ReadBodyAsJson(ctx, &utilityStock, false)
		utilityKey = utilityStock.StorehouseUtilityKey
	} else {
		// The key is the struct of a storehouse utility
		bodyreader.ReadBodyAsJson(ctx, &utility, false)
		utilityKey = utility.Key
	}
	_, errNotFound := dbVal.StorehouseUtilityDataSource.GetStorehouseUtilityByKey(utilityKey)
	if errNotFound != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("Un elemento con clave: %v no existe en almacen.", utilityKey),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseUtilityDeleted(ctx iris.Context) {
	var (
		httpRes models.Responser
		isError bool
	)
	utilityKey := ctx.Params().GetString("key")
	reqPath := ctx.Path()
	isUtilityDeleted, _ := dbVal.StorehouseUtilityDataSource.GetStorehouseUtilityByKey(utilityKey)
	var message string
	// if strings.Contains(reqPath, "pharmacy_inventory") && ctx.Request().Method == "POST" && isUtilityDeleted == 0 {
	// 	isError = true
	// 	message = fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", utilityKey)
	if strings.Contains(reqPath, "reactivate") && ctx.Request().Method == "PUT" && !isUtilityDeleted.DeletedAt.Valid {
		isError = true
		message = fmt.Sprintf("El elemento de almacen con clave: %v no se reactivó debido a que no ha sido eliminado antes.", utilityKey)
	}
	if isError {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    message,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) FindStorehouseStockById(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	storehouseStockId := ctx.Params().GetString("id")
	_, err := dbVal.StorehouseStocksDataSource.GetStorehouseStockById(storehouseStockId)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El stock con id: %v no existe.", storehouseStockId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseStockUsed(ctx iris.Context) {
	var (
		httpRes models.Responser
		action  string
	)
	storehouseStockId := ctx.Params().GetString("id")
	isStockUsed := dbVal.StorehouseStocksDataSource.IsStockUsed(storehouseStockId)
	if isStockUsed {
		if ctx.Method() == "PUT" {
			action = "actualizar"
		} else {
			action = "eliminar"
		}
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("No se puede %v un stock de almacen cuando ha sido utilizado.", action),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseStockWithLotNumber(ctx iris.Context) {
	var (
		storehouseStockId             string
		storehouseUtilityStockRequest models.StorehouseStock
		httpRes                       models.Responser
	)
	if ctx.Method() == "POST" {
		bodyreader.ReadBodyAsJson(ctx, &storehouseUtilityStockRequest, false)
	} else {
		storehouseStockId = ctx.Params().GetStringDefault("id", "")
		fmt.Println("id", storehouseStockId)
		bodyreader.ReadBodyAsJson(ctx, &storehouseUtilityStockRequest, false)
	}
	fmt.Println("lote", storehouseUtilityStockRequest.LotNumber)
	if storehouseUtilityStockRequest.LotNumber != "" {
		isFixedAsset := dbVal.StorehouseStocksDataSource.IsStorehouseStockWithLotNumber(storehouseUtilityStockRequest.LotNumber, storehouseStockId)

		fmt.Println("is", isFixedAsset)
		if isFixedAsset {
			httpRes = models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    fmt.Sprintf("Ya existe un elemento de almacén con el numero de lote: %v.", storehouseUtilityStockRequest.LotNumber),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseStockWithCatalogNumber(ctx iris.Context) {
	var (
		storehouseStockId             string
		storehouseUtilityStockRequest models.StorehouseStock
		httpRes                       models.Responser
	)
	if ctx.Method() == "POST" {
		bodyreader.ReadBodyAsJson(ctx, &storehouseUtilityStockRequest, false)
	} else {
		storehouseStockId = ctx.Params().GetStringDefault("id", "")
		bodyreader.ReadBodyAsJson(ctx, &storehouseUtilityStockRequest, false)
	}
	fmt.Println("catalogo", storehouseUtilityStockRequest.CatalogNumber, storehouseStockId)
	if storehouseUtilityStockRequest.CatalogNumber != "" {
		isFixedAsset := dbVal.StorehouseStocksDataSource.IsStorehouseStockWithCatalogNumber(storehouseUtilityStockRequest.CatalogNumber, storehouseStockId)
		fmt.Println("is", isFixedAsset)
		if isFixedAsset {
			httpRes = models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    fmt.Sprintf("Ya existe un elemento de almacén con el numero de catálogo: %v.", storehouseUtilityStockRequest.CatalogNumber),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) CanUtilityInfoBeChanged(ctx iris.Context) {
	var (
		httpRes              models.Responser
		storehouseUtilityReq models.StorehouseUtility
	)
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	bodyreader.ReadBodyAsJson(ctx, &storehouseUtilityReq, false)
	isStockUsed := dbVal.StorehouseUtilityDataSource.IsStockFromUtilityUsed(storehouseUtilityKey)
	storehouseUtility, _ := dbVal.StorehouseUtilityDataSource.GetStorehouseUtilityByKey(storehouseUtilityKey)
	isPresentationInfoChanged := storehouseUtilityReq.CategoryId != storehouseUtility.CategoryId || storehouseUtilityReq.PresentationId != storehouseUtility.PresentationId || storehouseUtilityReq.UnitId != storehouseUtility.UnitId
	isQuantityPerUnitChaged := storehouseUtilityReq.QuantityPerUnit != storehouseUtility.QuantityPerUnit
	if isStockUsed {
		if isPresentationInfoChanged || isQuantityPerUnitChaged {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    "La información de presentación o cantidad por unidad de un elemento de almacén no puede ser modificada una vez el elemento haya sido utilizado en una solicitud.",
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) IsStorehouseRequest(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	storehouseRequestId := ctx.Params().GetString("id")
	_, errNotFound := dbVal.StorehouseRequestsDataSource.GetStorehouseRequestById(storehouseRequestId)
	if errNotFound != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusForbidden,
			Message:    fmt.Sprintf("La solicitud con id: %v no existe.", storehouseRequestId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsSameRequestCreator(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	storehouseRequestId := ctx.Params().GetString("id")
	userId := ctx.Values().GetString("userId")
	isSameCreator := dbVal.StorehouseRequestsDataSource.IsSameRequestCreator(storehouseRequestId, userId)
	if !isSameCreator {
		httpRes = models.Responser{
			StatusCode: iris.StatusForbidden,
			Message:    "Solo el creador de la solicitud puede modificiarla/eliminarla",
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsRequestDeterminedStatus(statuses ...string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		var (
			httpRes models.Responser
		)
		storehouseRequestId := ctx.Params().GetString("id")
		IsDeterminedStatus := dbVal.StorehouseRequestsDataSource.IsRequestDeterminedStatus(storehouseRequestId, statuses)
		if !IsDeterminedStatus {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    fmt.Sprintf("No se pudó completar la acción, la solicitid no tiene un estado: %v", strings.Join(statuses, "/")),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
		ctx.Next()
	}
}

func (dbVal DbValidator) AreStorehouseRequestItemsValid(ctx iris.Context) {
	var (
		httpRes           models.Responser
		storehouseRequest models.StorehouseRequestDetailed
	)
	bodyreader.ReadBodyAsJson(ctx, &storehouseRequest, false)
	for _, requestUtility := range storehouseRequest.Utilities {
		utility, errNotFound := dbVal.StorehouseUtilityDataSource.GetStorehouseUtilityByKey(requestUtility.UtilityKey)
		if errNotFound != nil || utility.DeletedAt.Valid {
			httpRes = models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    fmt.Sprintf("El elemento de almacen con clave: %v no existe o se encuentra deshabilitado.", requestUtility.UtilityKey),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) AreFixedAssetsValidFromRequest(ctx iris.Context) {
	var (
		httpRes           models.Responser
		fixedAssetsRequet models.FixedAssetsRequestDetailed
	)
	bodyreader.ReadBodyAsJson(ctx, &fixedAssetsRequet, false)
	for _, fixedAssetItemRequest := range fixedAssetsRequet.FixedAssets {
		_, errNotFound := dbVal.FixedAssetsDataSource.GetFixedAssetByKey(fixedAssetItemRequest.FixedAsset.Key)
		if errNotFound == nil {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    fmt.Sprintf("El elemento de material de activo fijo con clave: %v ya se encuentra registrado.", fixedAssetItemRequest.FixedAsset.Key),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) FindFixedAssetByKey(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	fixedAssetKey := ctx.Params().GetStringDefault("key", "")
	_, errNotFound := dbVal.FixedAssetsDataSource.GetFixedAssetByKey(fixedAssetKey)
	if errNotFound != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El elemento de material de activo fijo con clave: %v no existe.", fixedAssetKey),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsFixedAssetWithKey(ctx iris.Context) {
	var (
		httpRes    models.Responser
		fixedAsset models.FixedAssetDetailed
	)
	bodyreader.ReadBodyAsJson(ctx, &fixedAsset, false)

	fixedAssetKey := ctx.Params().GetStringDefault("key", "")
	if fixedAsset.Key != fixedAssetKey {
		_, errNotFound := dbVal.FixedAssetsDataSource.GetFixedAssetByKey(fixedAsset.Key)
		if errNotFound == nil {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    fmt.Sprintf("No se pudo actualizar el elemento de material de activo fijo debido a que ya existe uno con clave: %v.", fixedAsset.Key),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) FindFixedAssetsRequestById(ctx iris.Context) {
	var httpRes models.Responser
	fixedAssetRequestId := ctx.Params().GetStringDefault("id", "")
	_, errNotFound := dbVal.FixedAssetsRequestsDataSource.GetFixedAssetsRequestById(fixedAssetRequestId)
	if errNotFound != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("La petición con id %v no existe.", fixedAssetRequestId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsFixedAssetWithSeries(ctx iris.Context) {
	var (
		fixedAssetKey      string
		fixedAsset         models.FixedAsset
		fixedAssetsRequest models.FixedAssetsRequestDetailed
		httpRes            models.Responser
	)
	if ctx.Method() == "POST" && strings.Contains(ctx.FullRequestURI(), "fixed_assets_requests") {
		bodyreader.ReadBodyAsJson(ctx, &fixedAssetsRequest, false)
	} else {
		fixedAssetKey = ctx.Params().GetStringDefault("key", "")
		bodyreader.ReadBodyAsJson(ctx, &fixedAsset, false)
		fixedAssetsRequest.FixedAssets = append(fixedAssetsRequest.FixedAssets, models.FixedAssetsItemsRequests{FixedAsset: fixedAsset})
	}
	for _, fixedAssetItemFromReq := range fixedAssetsRequest.FixedAssets {
		if fixedAssetItemFromReq.FixedAsset.Series != "" {
			isFixedAsset := dbVal.FixedAssetsDataSource.IsFixedAssetWithSeries(fixedAssetItemFromReq.FixedAsset.Series, fixedAssetKey)
			if isFixedAsset {
				httpRes = models.Responser{
					StatusCode: iris.StatusNotFound,
					Message:    fmt.Sprintf("Ya existe con material de activo fijo con el numero de serie: %v", fixedAssetItemFromReq.FixedAsset.Series),
				}
				helpers.PrepareAndSendMessageResponse(ctx, httpRes)
				return
			}
		}
	}
	ctx.Next()
}
