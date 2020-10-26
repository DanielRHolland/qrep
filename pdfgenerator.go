package main

import (
	"bytes"
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"log"
	"os"
	"time"
)

func generateQrsPdf(items []trackedItem) bytes.Buffer {
	begin := time.Now()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	for i := 0; i < len(items); i = i + 3 {
		m.Row(40, func() {
			for j := 0; j < 3; j++ {
				if index := i + j; index < len(items) {
					item := items[index]
					log.Println(index, item)
					m.Col(2, func() {

						m.QrCode(base+"/"+item.Id, props.Rect{
							Left:    5,
							Top:     5,
							Center:  false,
							Percent: 100,
						})
					})
					m.Col(2, func() {
						m.Text(item.Name, props.Text{
							Top:   5,
                                                        VerticalPadding: 5.0,
                                                        Align: consts.Center,
						})
					})

				}
			}
		})
	}

	//	err := m.OutputFileAndClose("example.pdf")
	pdfdata, err := m.Output()
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
	return pdfdata
}
