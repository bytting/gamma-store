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
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
    c.JSON(200, "pong")
}

func main() {
    db, err := sql.Open("postgres", "user=numsys dbname=gs sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    router := gin.Default()
    router.GET("/ping", ping)
    router.Run(":80")
}
