package mqttHealthcheck

import (
	"fmt"
	"net"
	"time"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

func MqttBroker(attrs *things.MqttBroker) error {
	// Test a TCP connection to the MQTT broker
	address := fmt.Sprintf("%s:%d", attrs.Host, attrs.Port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to MQTT broker at %s: %v", address, err)
	}
	defer conn.Close()

	return nil
}

func MqttClient(m *things.MqttClient) error {
	return fmt.Errorf("not implemented - perform a healthcheck for the MQTT broker instead")
}
