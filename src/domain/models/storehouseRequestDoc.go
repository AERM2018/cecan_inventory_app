package models

import (
	"cecan_inventory/domain/common"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type StorehouseRequestDoc struct {
	StorehouseRequest StorehouseRequestDetailed
	CreatedAt         time.Time
}

func (doc *StorehouseRequestDoc) CreateDoc() (string, error) {
	doc.CreatedAt = time.Now()
	cwd, _ := os.Getwd()
	outputFilePath := path.Join(cwd, "domain", "pdfs", fmt.Sprintf("%v_%v.pdf", "storehouseRequest", doc.CreatedAt.UnixMicro()))
	mto := pdf.NewMaroto(consts.Landscape, consts.Letter)
	common.InsertPdfHeader(mto)
	doc._addRequestInfo(mto)
	doc._addUtilitiesList(mto)
	err := mto.OutputFileAndClose(outputFilePath)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}
	return outputFilePath, nil
}

func (doc StorehouseRequestDoc) _addRequestInfo(m pdf.Maroto) {
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Solicitud de almacén", props.Text{
				Size:  18,
				Align: consts.Center,
				Style: consts.Bold,
			})
		})
	})
	// m.SetBorder(true)
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Folio: ", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
		m.Col(5, func() {
			m.Text(fmt.Sprintf("%v", doc.StorehouseRequest.Folio), props.Text{
				Size:  12,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text("Estatus: ", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
		m.Col(4, func() {
			m.Text(doc.StorehouseRequest.Status.Name, props.Text{
				Size:  12,
				Align: consts.Left,
			})
		})
	})
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Solicitante: ", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
		m.Col(5, func() {
			m.Text(doc.StorehouseRequest.User.FullName, props.Text{
				Size:  12,
				Align: consts.Left,
			})
		})
	})
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Fecha de creación: ", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
		m.Col(5, func() {
			m.Text(doc.StorehouseRequest.CreatedAt.Format(time.RFC822), props.Text{
				Size:  12,
				Align: consts.Left,
			})
		})
	})
}

func (doc StorehouseRequestDoc) _addUtilitiesList(m pdf.Maroto) {
	headers := []string{"Clave", "Nombre generico", "Categoría", "Cantidad solicitada", "Cantidad suminstrada"}
	gridSizes := []uint{2, 3, 2, 2, 3}
	contents := make([][]string, 0)
	for _, utility := range doc.StorehouseRequest.Utilities {
		utilityuInfo := []string{
			utility.StorehouseUtilty.Key,
			utility.StorehouseUtilty.GenericName,
			utility.StorehouseUtilty.Category.Name,
			fmt.Sprintf("%v %v", utility.Pieces, utility.StorehouseUtilty.Unit.Name),
			fmt.Sprintf("%v %v", utility.PiecesSupplied, utility.StorehouseUtilty.Unit.Name)}
		contents = append(contents, utilityuInfo)
	}
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.TableList(headers, contents, props.TableList{
		HeaderContentSpace:     0,
		VerticalContentPadding: 1.5,
		HeaderProp:             props.TableListContent{Size: 12, Style: consts.Bold, GridSizes: gridSizes},
		ContentProp:            props.TableListContent{Size: 10, GridSizes: gridSizes},
	})
	m.Line(15, props.Line{Color: color.NewWhite()})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Barcode(doc.StorehouseRequest.Id.String(), props.Barcode{Proportion: props.Proportion{Width: 80}, Center: true})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(doc.StorehouseRequest.Id.String(), props.Text{Align: consts.Center})
		})
	})
}
