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
	"time"

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

func apiSyncSession(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionName := c.Param("session_name")

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			abortApiRequest(c, http.StatusBadRequest, err)
			return
		}

		sync := new(Sync)
		if err := json.Unmarshal(body, sync); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		if len(sync.SessionIndices) > 60 {
			abortApiRequest(c, http.StatusRequestEntityTooLarge, err)
			return
		}

		spectrums, err := dbSelectSessionSync(db, sessionName, sync)
		if err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, spectrums)
	}
}

func apiAddSpectrum(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			abortApiRequest(c, http.StatusBadRequest, err)
			return
		}

		spec := new(Spectrum)
		if err := json.Unmarshal(body, spec); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		if err := dbInsertSpectrum(db, spec); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, makeApiResponseMessage("Spectrum inserted"))
	}
}

func apiGetSpectrums(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionName := c.Param("session_name")
		strDateBegin := c.Param("date_begin")
		strDateEnd := c.Param("date_end")

		const dateFormat string = "20060102_150405"

		if len(strDateBegin) == 0 {
			strDateBegin = "19000101_000000"
		}

		if len(strDateEnd) == 0 {
			strDateEnd = time.Now().Format(dateFormat)
		}

		dateBegin, err := time.Parse(dateFormat, strDateBegin)
		dateEnd, err := time.Parse(dateFormat, strDateEnd)

		spectrums, err := dbSelectSpectrums(db, sessionName, dateBegin, dateEnd)
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

	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/get-sessions", apiGetSessions(db))
	r.POST("/sync-session/:session_name", apiSyncSession(db))
	r.POST("/add-spectrum", apiAddSpectrum(db))
	r.GET("/get-spectrums/:session_name", apiGetSpectrums(db))
	r.GET("/get-spectrums/:session_name/:date_begin", apiGetSpectrums(db))
	r.GET("/get-spectrums/:session_name/:date_begin/:date_end", apiGetSpectrums(db))
	r.Run(":80")
}
