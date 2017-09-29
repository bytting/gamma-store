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
    "time"
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
    Doserate float64 `json:"doserate"`
}

func addSpectrum(c *gin.Context) {

    db, ok := c.Keys["db"].(*sql.DB)
    if !ok {
        fmt.Println("No database instance in context")
        return
    }

    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        fmt.Print(err)
        return
    }

    s := new(Spectrum)
    err = json.Unmarshal(body, s)
    if err != nil {
        fmt.Print(err)
        return
    }

    const dateFormat string = "2006-01-02T15:04:05.999Z"
    dt, err := time.Parse(dateFormat, s.StartTime)
    if err != nil {
        fmt.Print(err)
        return
    }

    _, err = db.Exec(`insert into spectrum (session_name, session_index,
    start_time, latitude, latitude_error, longitude, longitude_error,
    altitude, altitude_error, track, track_error, speed, speed_error,
    climb, climb_error, livetime, realtime, total_count, num_channels,
    channels, doserate) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20, $21)`,
    s.SessionName, s.SessionIndex, dt,
    s.Latitude, s.LatitudeError, s.Longitude, s.LongitudeError,
    s.Altitude, s.AltitudeError, s.Track, s.TrackError, s.Speed, s.SpeedError,
    s.Climb, s.ClimbError, s.Livetime, s.Realtime, s.TotalCount,
    s.NumChannels, s.Channels, s.Doserate)
    if err != nil {
        fmt.Print(err)
        return
    }

    c.JSON(200, "Spectrum inserted")
}

func getSpectrums(c *gin.Context) {

    c.JSON(200, "get-spectrums")
}

func main() {

    db, err := sql.Open("postgres", "host=localhost user=numsys dbname=gs sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
        c.Next()
    })

    r.POST("/add-spectrum", addSpectrum)
    r.GET("/get-spectrums", getSpectrums)
    r.Run(":80")
}

