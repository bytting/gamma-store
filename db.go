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
	"time"
)

const dbDateFormat string = "2006-01-02T15:04:05.999Z"

func dbConnectionString(hostname, username, dbname string) string {

	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", hostname, username, dbname)
}

func dbSelectSessions(db *sql.DB) ([]string, error) {

	rows, err := db.Query("select distinct session_name from spectrum")
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

func dbInsertSpectrum(db *sql.DB, s *Spectrum) error {

	dateTime, err := time.Parse(dbDateFormat, s.StartTime)
	if err != nil {
		return err
	}

	const sql_insert_spectrum = `
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

func dbSelectSpectrums(db *sql.DB, sessionName string, dateBegin, dateEnd time.Time) ([]Spectrum, error) {

	rows, err := db.Query("select * from spectrum where session_name = $1 and start_time between $2 and $3 order by start_time", sessionName, dateBegin, dateEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	spectrums := make([]Spectrum, 0)
	var s Spectrum
	for rows.Next() {
		var id int
		var dateTime time.Time
		if err := rows.Scan(&id, &s.SessionName, &s.SessionIndex, &dateTime, &s.Latitude, &s.Longitude, &s.Altitude, &s.Track, &s.Speed, &s.Climb, &s.Livetime, &s.Realtime, &s.NumChannels, &s.Channels, &s.Doserate); err != nil {
			return nil, err
		}
		s.StartTime = dateTime.Format(dbDateFormat)
		spectrums = append(spectrums, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spectrums, nil
}
