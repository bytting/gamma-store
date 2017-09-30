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

type Spectrum struct {
	SessionName  string  `json:"session_name"`
	SessionIndex int     `json:"session_index"`
	StartTime    string  `json:"start_time"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Altitude     float64 `json:"altitude"`
	Track        float64 `json:"track"`
	Speed        float64 `json:"speed"`
	Climb        float64 `json:"climb"`
	Livetime     float64 `json:"livetime"`
	Realtime     float64 `json:"realtime"`
	NumChannels  int     `json:"num_channels"`
	Channels     string  `json:"channels"`
	Doserate     float64 `json:"doserate"`
}
