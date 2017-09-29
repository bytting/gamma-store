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

const sql_connection_string = "host=localhost user=numsys dbname=gs sslmode=disable"

const sql_insert_spectrum = ` insert into spectrum (
    session_name, session_index, start_time,
    latitude, latitude_error, longitude,
    longitude_error, altitude, altitude_error,
    track, track_error, speed,
    speed_error, climb, climb_error,
    livetime, realtime, total_count,
    num_channels, channels, doserate
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`
