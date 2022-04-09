package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	response, err := http.Get("https://munroapi.herokuapp.com/munros")
	if err != nil {
		fmt.Print(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err)
	}

	munros, err := UnmarshalMunros(data)
	if err != nil {
		fmt.Print(err)
	}

	a := app.New()
	w := a.NewWindow("Munros")
	w.Resize(fyne.NewSize(800, 800))

	munroList := widget.NewList(
		func() int { return len(munros) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(munros[id].Name)
		})

	munroText := widget.NewLabel("Select a Munro for more information.")

	munroList.OnSelected = func(id widget.ListItemID) {
		munroDetail := fmt.Sprintf("Name: %s \nHeight: %d \nRegion: %s \nMeaning: %s",
			munros[id].Name, munros[id].Height, munros[id].Region, munros[id].Meaning)
		munroText.SetText(munroDetail)
	}

	split := container.NewHSplit(
		munroList,
		container.NewMax(munroText),
	)

	split.Offset = 0.4

	w.SetContent(split)

	w.ShowAndRun()

}

type Munros []Munro

func UnmarshalMunros(data []byte) (Munros, error) {
	var r Munros
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Munros) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Munro struct {
	Name             string         `json:"name"`
	Height           int64          `json:"height"`
	GridrefLetters   GridrefLetters `json:"gridref_letters"`
	GridrefEastings  string         `json:"gridref_eastings"`
	GridrefNorthings string         `json:"gridref_northings"`
	LatlngLat        float64        `json:"latlng_lat"`
	LatlngLng        float64        `json:"latlng_lng"`
	Smcid            string         `json:"smcid"`
	MetofficeLOCID   string         `json:"metoffice_loc_id"`
	Region           string         `json:"region"`
	Meaning          string         `json:"meaning"`
}

type GridrefLetters string

const (
	Nc GridrefLetters = "NC"
	Ng GridrefLetters = "NG"
	Nh GridrefLetters = "NH"
	Nj GridrefLetters = "NJ"
	Nm GridrefLetters = "NM"
	Nn GridrefLetters = "NN"
	No GridrefLetters = "NO"
)
