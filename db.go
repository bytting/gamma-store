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
    "fmt"
    "strings"
    "time"
)

const dbDateFormat string = "2006-01-02T15:04:05.999Z"

func dbConnectionString(hostname, username, dbname string) string {

    return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", hostname, username, dbname)
}

func dbValidateCredentials(db *sql.DB, user, pass string) (bool, error) {

    rows, err := db.Query("select id from users where username = $1 and password = $2", user, pass)
    if err != nil {
        return false, err
    }
    defer rows.Close()

    if !rows.Next() {
        return false, nil
    }

    return true, nil
}

func dbSelectSessions(db *sql.DB) ([]string, error) {

    rows, err := db.Query("select distinct session_name from spectrum order by session_name desc")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    sessionNames := make([]string, 0)
    var sessionName string
    for rows.Next() {
        if err := rows.Scan(&sessionName); err != nil {
            return nil, err
        }
        sessionNames = append(sessionNames, sessionName)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return sessionNames, nil
}

func dbSelectSessionSync(db *sql.DB, sessionName string, sync *Sync) ([]Spectrum, error) {

    indexList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sync.SessionIndices)), ","), "[]")
    query := fmt.Sprintf("select * from spectrum where session_name = $1 and (session_index in (%s) or session_index > $2)", indexList)

    rows, err := db.Query(query, sessionName, sync.LastIndex)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    spectrums := make([]Spectrum, 0)
    var spec Spectrum
    for rows.Next() {
        var id int
        var dateTime time.Time
        err = rows.Scan(
            &id,
            &spec.SessionName,
            &spec.SessionIndex,
            &dateTime,
            &spec.Latitude,
            &spec.Longitude,
            &spec.Altitude,
            &spec.Track,
            &spec.Speed,
            &spec.Climb,
            &spec.Livetime,
            &spec.Realtime,
            &spec.NumChannels,
            &spec.Channels,
            &spec.Doserate)
        if err != nil {
            return nil, err
        }

        spec.StartTime = dateTime.Format(dbDateFormat)
        spectrums = append(spectrums, spec)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return spectrums, nil
}

func dbAddSpectrum(db *sql.DB, s *Spectrum) error {

    sql_insert_spectrum := "select id from spectrum where session_name = $1 and session_index = $2"

    rows, err := db.Query(sql_insert_spectrum,
        s.SessionName,
        s.SessionIndex)
    if err != nil {
        return err
    }
    defer rows.Close()

    if !rows.Next() {
        if rows.Err() != nil {
            return rows.Err()
        }
        return dbInsertSpectrum(db, s)

    } else {
        return dbUpdateSpectrum(db, s)
    }
}

func dbInsertSpectrum(db *sql.DB, s *Spectrum) error {

    dateTime, err := time.Parse(dbDateFormat, s.StartTime)
    if err != nil {
        return err
    }

    sql_insert_spectrum := `
        insert into spectrum (
        session_name,
        session_index,
        start_time,
        latitude,
        longitude,
        altitude,
        track,
        speed,
        climb,
        livetime,
        realtime,
        num_channels,
        channels,
        doserate
        ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

    _, err = db.Exec(sql_insert_spectrum,
        s.SessionName,
        s.SessionIndex,
        dateTime,
        s.Latitude,
        s.Longitude,
        s.Altitude,
        s.Track,
        s.Speed,
        s.Climb,
        s.Livetime,
        s.Realtime,
        s.NumChannels,
        s.Channels,
        s.Doserate)

    return err
}

func dbUpdateSpectrum(db *sql.DB, s *Spectrum) error {

    dateTime, err := time.Parse(dbDateFormat, s.StartTime)
    if err != nil {
        return err
    }

    sql_update_spectrum := `
        update spectrum set
        start_time = $1,
        latitude = $2,
        longitude = $3,
        altitude = $4,
        track = $5,
        speed = $6,
        climb = $7,
        livetime = $8,
        realtime = $9,
        num_channels = $10,
        channels = $11,
        doserate = $12
        where session_name = $13 and session_index = $14`

    _, err = db.Exec(sql_update_spectrum,
        dateTime,
        s.Latitude,
        s.Longitude,
        s.Altitude,
        s.Track,
        s.Speed,
        s.Climb,
        s.Livetime,
        s.Realtime,
        s.NumChannels,
        s.Channels,
        s.Doserate,
        s.SessionName,
        s.SessionIndex)

    return err
}

func dbSelectSpectrums(db *sql.DB, sessionName string, dateBegin, dateEnd time.Time) ([]Spectrum, error) {

    sql_select_spectrum := "select * from spectrum where session_name = $1 and start_time between $2 and $3 order by start_time"

    rows, err := db.Query(sql_select_spectrum,
        sessionName,
        dateBegin,
        dateEnd)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    spectrums := make([]Spectrum, 0)
    var spec Spectrum
    for rows.Next() {
        var id int
        var dateTime time.Time
        err = rows.Scan(
            &id,
            &spec.SessionName,
            &spec.SessionIndex,
            &dateTime,
            &spec.Latitude,
            &spec.Longitude,
            &spec.Altitude,
            &spec.Track,
            &spec.Speed,
            &spec.Climb,
            &spec.Livetime,
            &spec.Realtime,
            &spec.NumChannels,
            &spec.Channels,
            &spec.Doserate)
        if err != nil {
            return nil, err
        }

        spec.StartTime = dateTime.Format(dbDateFormat)
        spectrums = append(spectrums, spec)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return spectrums, nil
}

