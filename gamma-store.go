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
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func abortApiRequest(c *gin.Context, status int, err error) {

	c.AbortWithStatus(status)
	log.Print(err)
}

func makeApiResponseMessage(msg string) gin.H {

	return gin.H{"message": msg}
}

func apiGetSessions(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		sessions, err := dbSelectSessions(db)
		if err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

func apiAddSpectrum(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			abortApiRequest(c, http.StatusBadRequest, err)
			return
		}

		s := new(Spectrum)
		if err := json.Unmarshal(body, s); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		if err := dbInsertSpectrum(db, s); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, makeApiResponseMessage("Spectrum inserted"))
	}
}

func apiGetSpectrums(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		session := c.Param("session")

		spectrums, err := dbSelectSpectrums(db, session)
		if err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, spectrums)
	}
}

func main() {

	db, err := sql.Open("postgres", dbConnectionString("localhost", "numsys", "gs"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/get-sessions", apiGetSessions(db))
	r.POST("/add-spectrum", apiAddSpectrum(db))
	r.GET("/get-spectrums/:session", apiGetSpectrums(db))
	r.Run(":80")
}
