/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"context"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MiHomeMode       uint
	MiHomeProduct    byte
	MiHomeValveState byte
	MiHomePowerMode  byte
)

////////////////////////////////////////////////////////////////////////////////
// ENER314 AND ENER314RT

type ENER314 interface {
	gopi.Driver

	// Send on signal - when no sockets specified then
	// sends to all sockets
	On(sockets ...uint) error

	// Send off signal - when no sockets specified then
	// sends to all sockets
	Off(sockets ...uint) error
}

type ENER314RT interface {

	// Receive payloads with radio until context deadline exceeded or cancel,
	// this blocks sending
	Receive(ctx context.Context, mode MiHomeMode, payload chan<- []byte) error

	// Send a raw payload with radio
	Send(payload []byte, repeat uint, mode MiHomeMode) error

	// Measure device temperature
	MeasureTemperature(offset float32) (float32, error)

	// ResetRadio device
	ResetRadio() error
}

////////////////////////////////////////////////////////////////////////////////
// MIHOME

type MiHome interface {
	gopi.Publisher

	// Reset the device
	Reset() error

	// Add a wire protocol which encodes/decodes messages
	AddProto(Proto) error

	// Return registered protocols
	Protos() []Proto

	// Measure Device Temperature
	MeasureTemperature() (float32, error)

	// Request Switch state for both monitor and control devices
	RequestSwitchOn(MiHomeProduct, uint32) error
	RequestSwitchOff(MiHomeProduct, uint32) error

	// Send a join message after a report is received
	SendJoin(MiHomeProduct, uint32) error

	// Note the eTRV messages below should be sent very shortly
	// after the temperature report is provided as that's when the
	// eTRV is awake to respond to the messages
	RequestIdentify(MiHomeProduct, uint32) error
	RequestDiagnostics(MiHomeProduct, uint32) error
	RequestExercise(MiHomeProduct, uint32) error
	RequestBatteryLevel(MiHomeProduct, uint32) error
	RequestTargetTemperature(MiHomeProduct, uint32, float64) error
	RequestReportInterval(MiHomeProduct, uint32, time.Duration) error
	RequestValveState(MiHomeProduct, uint32, MiHomeValveState) error
	RequestLowPowerMode(MiHomeProduct, uint32, bool) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MIHOME_MODE_NONE    MiHomeMode = iota
	MIHOME_MODE_MONITOR            // FSK
	MIHOME_MODE_CONTROL            // OOK
	MIHOME_MODE_MAX     = MIHOME_MODE_CONTROL
)

const (
	// Monitor Products (FSK)
	MIHOME_PRODUCT_NONE    MiHomeProduct = 0x00
	MIHOME_PRODUCT_MIHO004 MiHomeProduct = 0x01 // Adaptor Monitor
	MIHOME_PRODUCT_MIHO005 MiHomeProduct = 0x02 // Adaptor Plus
	MIHOME_PRODUCT_MIHO013 MiHomeProduct = 0x03 // eTRV
	MIHOME_PRODUCT_MIHO006 MiHomeProduct = 0x05 // House Monitor
	MIHOME_PRODUCT_MIHO032 MiHomeProduct = 0x0C // Motion sensor
	MIHOME_PRODUCT_MIHO033 MiHomeProduct = 0x0D // Door sensor

	// Control Products (OOK)
	MIHOME_PRODUCT_CONTROL_ALL   MiHomeProduct = 0xF0 // OOK Switch All
	MIHOME_PRODUCT_CONTROL_ONE   MiHomeProduct = 0xF1 // OOK Switch 1
	MIHOME_PRODUCT_CONTROL_TWO   MiHomeProduct = 0xF2 // OOK Switch 2
	MIHOME_PRODUCT_CONTROL_THREE MiHomeProduct = 0xF3 // OOK Switch 3
	MIHOME_PRODUCT_CONTROL_FOUR  MiHomeProduct = 0xF4 // OOK Switch 4
)

const (
	MIHOME_VALVE_STATE_OPEN   MiHomeValveState = 0x00 // Valve fully open
	MIHOME_VALVE_STATE_CLOSED MiHomeValveState = 0x01 // Valve fully closed
	MIHOME_VALVE_STATE_NORMAL MiHomeValveState = 0x02 // Valve in normal state
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Mode returns the mode for a product
func (p MiHomeProduct) Mode() MiHomeMode {
	switch p {
	case MIHOME_PRODUCT_MIHO004:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_MIHO005:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_MIHO013:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_MIHO006:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_MIHO032:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_MIHO033:
		return MIHOME_MODE_MONITOR
	case MIHOME_PRODUCT_CONTROL_ALL:
		return MIHOME_MODE_CONTROL
	case MIHOME_PRODUCT_CONTROL_ONE:
		return MIHOME_MODE_CONTROL
	case MIHOME_PRODUCT_CONTROL_TWO:
		return MIHOME_MODE_CONTROL
	case MIHOME_PRODUCT_CONTROL_THREE:
		return MIHOME_MODE_CONTROL
	case MIHOME_PRODUCT_CONTROL_FOUR:
		return MIHOME_MODE_CONTROL
	default:
		return MIHOME_MODE_NONE
	}
}

// Socket returns the socket number of a control product
func (p MiHomeProduct) Socket() uint {
	switch p {
	case MIHOME_PRODUCT_CONTROL_ALL:
		return 0
	case MIHOME_PRODUCT_CONTROL_ONE:
		return 1
	case MIHOME_PRODUCT_CONTROL_TWO:
		return 2
	case MIHOME_PRODUCT_CONTROL_THREE:
		return 3
	case MIHOME_PRODUCT_CONTROL_FOUR:
		return 4
	default:
		// Return 0 otherwise
		return 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m MiHomeMode) String() string {
	switch m {
	case MIHOME_MODE_NONE:
		return "MIHOME_MODE_NONE"
	case MIHOME_MODE_MONITOR:
		return "MIHOME_MODE_MONITOR"
	case MIHOME_MODE_CONTROL:
		return "MIHOME_MODE_CONTROL"
	default:
		return "[?? Invalid MiHomeMode value]"
	}
}

func (p MiHomeProduct) String() string {
	switch p {
	case MIHOME_PRODUCT_NONE:
		return "MIHOME_PRODUCT_NONE"
	case MIHOME_PRODUCT_MIHO004:
		return "MIHOME_PRODUCT_MIHO004"
	case MIHOME_PRODUCT_MIHO005:
		return "MIHOME_PRODUCT_MIHO005"
	case MIHOME_PRODUCT_MIHO013:
		return "MIHOME_PRODUCT_MIHO013"
	case MIHOME_PRODUCT_MIHO006:
		return "MIHOME_PRODUCT_MIHO006"
	case MIHOME_PRODUCT_MIHO032:
		return "MIHOME_PRODUCT_MIHO032"
	case MIHOME_PRODUCT_MIHO033:
		return "MIHOME_PRODUCT_MIHO033"
	case MIHOME_PRODUCT_CONTROL_ALL:
		return "MIHOME_PRODUCT_CONTROL_ALL"
	case MIHOME_PRODUCT_CONTROL_ONE:
		return "MIHOME_PRODUCT_CONTROL_ONE"
	case MIHOME_PRODUCT_CONTROL_TWO:
		return "MIHOME_PRODUCT_CONTROL_TWO"
	case MIHOME_PRODUCT_CONTROL_THREE:
		return "MIHOME_PRODUCT_CONTROL_THREE"
	case MIHOME_PRODUCT_CONTROL_FOUR:
		return "MIHOME_PRODUCT_CONTROL_FOUR"
	default:
		return "[?? Invalid MiHomeProduct value]"
	}
}

func (s MiHomeValveState) String() string {
	switch s {
	case MIHOME_VALVE_STATE_OPEN:
		return "MIHOME_VALVE_STATE_OPEN"
	case MIHOME_VALVE_STATE_CLOSED:
		return "MIHOME_VALVE_STATE_CLOSED"
	case MIHOME_VALVE_STATE_NORMAL:
		return "MIHOME_VALVE_STATE_NORMAL"
	default:
		return "[?? Invalid MiHomeValveState value]"
	}
}
