package models

type Click struct {
	link      Link
	completed bool
	country   string // Country code eg. en
	language  string // Language eg. en-US
	browser   string // Browser name eg. chrome
	device_os string // Device os eg. Android, IOS, Mac, Windows
	ip        int    // Access ip
	click_ts  int    // Click timestamp
}
