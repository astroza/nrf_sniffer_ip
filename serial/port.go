/* (c) 2018 felipe@astroza.cl - See LICENSE
 */
package serial

import (
	"os"

	"github.com/astroza/nrf_sniffer_ip/common"
	"github.com/pkg/term"
)

const defaultSerialDevice = "/dev/ttyACM0"
const serialSpeed = 460800

func OpenPort() (*common.ReadWriterChannel, error) {
	var serialDevice string
	if serialDevice = os.Getenv("SERIAL_DEVICE"); serialDevice == "" {
		serialDevice = defaultSerialDevice
	}
	dev, err := term.Open(serialDevice, term.Speed(serialSpeed), term.RawMode, term.FlowControl(term.NONE))
	if err != nil {
		return nil, err
	}
	r := common.ReadWriterChannel{}
	r.Init(dev)
	return &r, nil
}
