# Device Registry service

## Usage
All responses have the form

```json
{
	"data":"mixed data",
	"message":"description of what happened"
}
``` 

Subsequent 

### List all devices

**defenition**

`GET /devices`

**Response**

- `200 OK` on success

```json
[
	{
		"identifier":"<device_identifier>",
		"name":"<device_name>",
		"device_type":"<type e.g. switch>",
		"controller_gateway":"<device_gateway>"
	}
]
```

### Add a new device

**defenition**

`POST /devices`

**Arguments**

`"identifier":string'`
`"name":string'`
`"type":string'`
`"controller_gateway":ip address'`

**Response**

- `201 Created` - on success

```json
{
	"identifier":"<device_identifier>",
	"name":"<device_name>",
	"device_type":"<type e.g. switch>",
	"controller_gateway":"<device_gateway>"
}
```

### Lookup device details

**defenition**

`GET /device/<identifier>`

**Response**

- `404 Not Found` if it does not exist
- `200 OK` on success

```json
{
	"identifier":"<device_identifier>",
	"name":"<device_name>",
	"device_type":"<type e.g. switch>",
	"controller_gateway":"<device_gateway>"
}
```

### Delete a device

**defenition**

`DELETE /device/<identifier>`

**Response**

- `404 Not Found` if it does not exist
- `204 No Content` on success


## Types

Bits that are turned on (Left to Right)
0 - Controller
1 - Switch
2 - Microphone
3 - Temperature
4 - Humidity
5 - Phototometer
6 - Wireless transmission: 0 - WiFi, 1 - Bluetooth
7 - Power source: 0 - Battery, 1 - PowerPlug