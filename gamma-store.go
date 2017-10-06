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
    "encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"
    "strings"

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

func checkCredentials(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

        items := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
        if len(items) != 2 || items[0] != "Basic" {
            log.Println("Malformed credentials")
            c.JSON(http.StatusBadRequest, makeApiResponseMessage("Malformed credentials"))
            c.Abort()
            return
        }

        data, _ := base64.StdEncoding.DecodeString(items[1])
        cred := strings.SplitN(string(data), ":", 2)
        if len(cred) != 2 {
            log.Println("Malformed credentials")
            c.JSON(http.StatusBadRequest, makeApiResponseMessage("Malformed credentials"))
            c.Abort()
            return
        }

        valid, err := dbValidateCredentials(db, cred[0], cred[1])
        if err != nil {
            log.Println("Credential validation failed")
            c.JSON(http.StatusInternalServerError, makeApiResponseMessage("Credential validation failed"))
            c.Abort()
            return
        }

        if !valid {
            log.Println("Invalid credentials")
            c.JSON(http.StatusUnauthorized, makeApiResponseMessage("Invalid credentials"))
            c.Abort()
            return
        }

        c.Next()
    }
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
			abortApiRequest(c, http.StatusBadRequest, err)
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

		if err := dbAddSpectrum(db, spec); err != nil {
			abortApiRequest(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, makeApiResponseMessage("Spectrum stored"))
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
		if err != nil {
			abortApiRequest(c, http.StatusBadRequest, err)
			return
		}

		dateEnd, err := time.Parse(dateFormat, strDateEnd)
		if err != nil {
			abortApiRequest(c, http.StatusBadRequest, err)
			return
		}

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

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
    root := router.Group("/")
    root.Use(checkCredentials(db))
    {
        root.GET("/get-sessions", apiGetSessions(db))
        root.POST("/sync-session/:session_name", apiSyncSession(db))
        root.POST("/add-spectrum", apiAddSpectrum(db))
        root.GET("/get-spectrums/:session_name", apiGetSpectrums(db))
        root.GET("/get-spectrums/:session_name/:date_begin", apiGetSpectrums(db))
        root.GET("/get-spectrums/:session_name/:date_begin/:date_end", apiGetSpectrums(db))
    }

	router.Run(":80")
}
