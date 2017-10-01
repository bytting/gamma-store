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
	"errors"
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

func apiGetSessions(c *gin.Context) {

	db, ok := c.Keys["db"].(*sql.DB)
	if !ok {
		abortApiRequest(c, http.StatusInternalServerError, errors.New("Invalid database handle in context"))
		return
	}

	sessions, err := dbSelectSessions(db)
	if err != nil {
		abortApiRequest(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func apiAddSpectrum(c *gin.Context) {

	db, ok := c.Keys["db"].(*sql.DB)
	if !ok {
		abortApiRequest(c, http.StatusInternalServerError, errors.New("Invalid database handle in context"))
		return
	}

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

func apiGetSpectrums(c *gin.Context) {

	db, ok := c.Keys["db"].(*sql.DB)
	if !ok {
		abortApiRequest(c, http.StatusInternalServerError, errors.New("Invalid database handle in context"))
		return
	}

	spectrums, err := dbSelectSpectrums(db)
	if err != nil {
		abortApiRequest(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, spectrums)
}

func main() {

	db, err := sql.Open("postgres", dbConnectionString("localhost", "numsys", "gs"))
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

	r.GET("/get-sessions", apiGetSessions)
	r.POST("/add-spectrum", apiAddSpectrum)
	r.GET("/get-spectrums", apiGetSpectrums)
	r.Run(":80")
}
