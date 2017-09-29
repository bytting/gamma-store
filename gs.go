//    gamma-store - Web service to store gamma spectrum data
//
//    Copyright (C) 2017  NRPA
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
//    Authors: Dag Robole,

package main

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"
)

type Spectrum struct {
    SessionName string `json:"session_name"`
    SessionIndex int `json:"session_index"`
    StartTime string `json:"start_time"`
    Latitude float64 `json:"latitude"`
    LatitudeError float64 `json:"latitude_error"`
    Longitude float64 `json:"longitude"`
    LongitudeError float64 `json:"longitude_error"`
    Altitude float64 `json:"altitude"`
    AltitudeError float64 `json:"altitude_error"`
    Track float64 `json:"track"`
    TrackError float64 `json:"track_error"`
    Speed float64 `json:"speed"`
    SpeedError float64 `json:"speed_error"`
    Climb float64 `json:"climb"`
    ClimbError float64 `json:"climb_error"`
    Livetime float64 `json:"livetime"`
    Realtime float64 `json:"realtime"`
    TotalCount int `json:"total_count"`
    NumChannels int `json:"num_channels"`
    Channels string `json:"channels"`
    Doserate string `json:"doserate"`
}

func addSpectrum(c *gin.Context) {

    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        fmt.Print(err)
        return
    }

    spectrum := new(Spectrum)
    err = json.Unmarshal(body, spectrum)
    if err != nil {
        fmt.Print(err)
        return
    }

    c.JSON(200, spectrum)
}

func getSpectrums(c *gin.Context) {
    c.JSON(200, "get-spectrums")
}

func main() {

    db, err := sql.Open("postgres", "user=numsys dbname=gs sslmode=disable")
    if err != nil {
        fmt.Print(err)
    }
    defer db.Close()

    r := gin.Default()
    r.POST("/add-spectrum", addSpectrum)
    r.GET("/get-spectrums", getSpectrums)
    r.Run(":80")
}

