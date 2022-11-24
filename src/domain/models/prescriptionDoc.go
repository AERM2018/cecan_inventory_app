package models

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type PrescriptionDoc struct {
	Prescription PrescriptionDetialed
	CreatedAt    time.Time
}

func (doc *PrescriptionDoc) CreateDoc() (string, error) {
	cwd, _ := os.Getwd()
	doc.CreatedAt = time.Now()
	mto := pdf.NewMaroto(consts.Portrait, consts.Letter)
	mto.SetPageMargins(12, 12, 12)
	doc._addHeader(mto)
	doc._addPrescriptionInfo(mto)
	doc._addMedicineList(mto)
	outputFilePath := path.Join(cwd, "domain", "pdfs", "prescription.pdf")
	err := mto.OutputFileAndClose(outputFilePath)
	if err != nil {
		return "", err
	}
	return outputFilePath, nil

}
func (doc PrescriptionDoc) _addHeader(m pdf.Maroto) {
	cwd, _ := os.Getwd()
	m.RegisterHeader(func() {
		m.Row(15, func() {
			m.Col(1, func() {
				m.FileImage(
					path.Join(cwd, "domain", "assets", "salud_dgo.png"),
					props.Rect{Percent: 100, Top: 1},
				)
			})
			m.Col(1, func() {
				m.FileImage(
					path.Join(cwd, "domain", "assets", "cecan.png"),
					props.Rect{Percent: 100},
				)
			})
			m.Col(10, func() {
				m.Text("Centro de cancerología del estado de Durango", props.Text{
					Size:  18,
					Top:   3,
					Left:  5,
					Style: consts.Bold,
					Align: consts.Left,
				})
			})
		})
	})

}

func (doc PrescriptionDoc) _addPrescriptionInfo(m pdf.Maroto) {
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Receta digital", props.Text{
				Size:  16,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.Row(10, func() {
		m.Col(1, func() {
			m.Text("Folio: ", props.Text{Style: consts.Bold, Size: 12})
		})
		m.Col(1, func() {
			m.Text(fmt.Sprintf("%v", doc.Prescription.Folio), props.Text{Align: consts.Left, Size: 12})
		})
	})
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Estatus: ", props.Text{Style: consts.Bold, Size: 12})
		})
		m.Col(2, func() {
			m.Text(fmt.Sprintf("%v", doc.Prescription.PrescriptionStatus.Name), props.Text{Align: consts.Left, Size: 12})
		})
	})
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Doctor que suscribe: ", props.Text{Style: consts.Bold, Size: 12})
		})
		m.Col(5, func() {
			m.Text(fmt.Sprintf("%v (%v)", strings.ToTitle(doc.Prescription.User.FullName), doc.Prescription.User.Id), props.Text{Align: consts.Left, Size: 12})
		})
	})
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Nombre del paciente: ", props.Text{Style: consts.Bold, Size: 12})
		})
		m.Col(3, func() {
			m.Text(doc.Prescription.PatientName, props.Text{Align: consts.Left, Size: 12})
		})
	})
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Fecha de creación: ", props.Text{Style: consts.Bold, Size: 12})
		})
		m.Col(3, func() {
			m.Text(doc.Prescription.CreatedAt.Format(time.RFC822), props.Text{Align: consts.Left, Size: 12})
		})
	})
	if doc.Prescription.PrescriptionStatus.Name != "Pendiente" {
		m.Row(10, func() {
			m.Col(3, func() {
				m.Text("Fecha de ultima suministración: ", props.Text{Style: consts.Bold, Size: 12})
			})
			m.Col(3, func() {
				m.Text(doc.Prescription.SuppliedAt.Format(time.RFC822), props.Text{Align: consts.Left, Size: 12})
			})
		})
	}
}

func (doc PrescriptionDoc) _addMedicineList(m pdf.Maroto) {
	headers := []string{"Clave", "Nombre", "Piezas recetadas", "Piezas suminstradas"}
	gridSizes := []uint{2, 4, 3, 3}
	contents := make([][]string, 0)
	for _, medicine := range doc.Prescription.Medicines {
		medicineInfo := []string{medicine.Medicine.Key, medicine.Medicine.Name, fmt.Sprintf("%v", medicine.Pieces), fmt.Sprintf("%v", medicine.PiecesSupplied)}
		contents = append(contents, medicineInfo)
	}
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.TableList(headers, contents, props.TableList{
		HeaderContentSpace:     0,
		VerticalContentPadding: 1.5,
		HeaderProp:             props.TableListContent{Size: 12, Style: consts.Bold, GridSizes: gridSizes},
		ContentProp:            props.TableListContent{Size: 10, GridSizes: gridSizes},
	})
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Instrucciones: \n%v", doc.Prescription.Instructions), props.Text{Style: consts.Bold})
		})
	})
	m.Line(10, props.Line{Color: color.NewWhite()})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Barcode(doc.Prescription.Id.String(), props.Barcode{Proportion: props.Proportion{Width: 80}, Center: true})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(doc.Prescription.Id.String(), props.Text{Align: consts.Center})
		})
	})
}
