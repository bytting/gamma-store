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
    "log"
    "io/ioutil"
    "encoding/json"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"
)

type DetectorData struct {
    TypeName string `json:"type_name"`
    SerialNumber string `json:"serialnumber"`
}

type Session struct {
    Name string `json:"name"`
    Comment string `json:"comment"`
    *DetectorData `json:"detector_data"`
}

func addSession(c *gin.Context) {

    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Print(err)
        return
    }

    session := new(Session)
    err = json.Unmarshal(body, session)
    if err != nil {
        log.Print(err)
        return
    }

    c.JSON(200, session)
}

func getSessions(c *gin.Context) {
    c.JSON(200, "get-sessions")
}

func addSpectrum(c *gin.Context) {
    c.JSON(200, "add-spectrum")
}

func getSpectrums(c *gin.Context) {
    c.JSON(200, "get-spectrums")
}

func main() {

    db, err := sql.Open("postgres", "user=numsys dbname=gs sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := gin.Default()
    r.POST("/add-session", addSession)
    r.GET("/get-sessions", getSessions)
    r.POST("/add-spectrum", addSpectrum)
    r.GET("/get-spectrums", getSpectrums)
    r.Run(":80")
}

