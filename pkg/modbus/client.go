package modbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grid-x/modbus"
)

const (
	// Maximum number of registers that can be read in a single request
	MaxRegistersPerRead = 125
)

// Client represents a MODBUS TCP client
type Client struct {
	client  modbus.Client
	handler *modbus.TCPClientHandler

	mutex       sync.RWMutex
	isConnected bool
}

// NewClient creates a new MODBUS TCP client
func NewClient(host string, port int, slaveID byte, timeout time.Duration) *Client {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", host, port))
	handler.SlaveID = slaveID
	handler.Timeout = timeout

	client := modbus.NewClient(handler)

	return &Client{
		client:  client,
		handler: handler,
	}
}

// Connect establishes connection to the MODBUS server
func (c *Client) Connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.handler.Connect(ctx)
	if err != nil {
		c.isConnected = false
		return err
	}
	c.isConnected = true
	return nil
}

// Disconnect closes the connection to the MODBUS server
func (c *Client) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.handler.Close()
	c.isConnected = false
	return err
}

// IsConnected returns the current connection status
func (c *Client) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// GetSlaveID returns the current slave ID
func (c *Client) GetSlaveID() byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.handler.SlaveID
}

// SetSlaveID sets the slave ID for subsequent operations
func (c *Client) SetSlaveID(slaveID byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handler.SlaveID = slaveID
}

// ReadHoldingRegisters reads holding registers from the MODBUS server
func (c *Client) ReadHoldingRegisters(ctx context.Context, address, quantity uint16) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return nil, fmt.Errorf("modbus client not connected")
	}

	data, err := c.client.ReadHoldingRegisters(ctx, address, quantity)
	if err != nil {
		c.handleConnectionError(err)
		return nil, err
	}
	return data, nil
}

// ReadInputRegisters reads input registers from the MODBUS server
func (c *Client) ReadInputRegisters(ctx context.Context, address, quantity uint16) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return nil, fmt.Errorf("modbus client not connected")
	}

	data, err := c.client.ReadInputRegisters(ctx, address, quantity)
	if err != nil {
		c.handleConnectionError(err)
		return nil, err
	}
	return data, nil
}

// ReadDiscreteInputs reads discrete inputs from the MODBUS server
func (c *Client) ReadDiscreteInputs(ctx context.Context, address, quantity uint16) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return nil, fmt.Errorf("modbus client not connected")
	}

	data, err := c.client.ReadDiscreteInputs(ctx, address, quantity)
	if err != nil {
		c.handleConnectionError(err)
		return nil, err
	}
	return data, nil
}

// ReadCoils reads coils from the MODBUS server
func (c *Client) ReadCoils(ctx context.Context, address, quantity uint16) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return nil, fmt.Errorf("modbus client not connected")
	}

	data, err := c.client.ReadCoils(ctx, address, quantity)
	if err != nil {
		c.handleConnectionError(err)
		return nil, err
	}
	return data, nil
}

// WriteSingleRegister writes a single register to the MODBUS server
func (c *Client) WriteSingleRegister(ctx context.Context, address, value uint16) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("modbus client not connected")
	}

	_, err := c.client.WriteSingleRegister(ctx, address, value)
	if err != nil {
		c.handleConnectionError(err)
		return err
	}
	return nil
}

// WriteMultipleRegisters writes multiple registers to the MODBUS server
func (c *Client) WriteMultipleRegisters(ctx context.Context, address uint16, values []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("modbus client not connected")
	}

	if len(values)%2 != 0 {
		return fmt.Errorf("values must be even number of bytes, got %d", len(values))
	}

	_, err := c.client.WriteMultipleRegisters(ctx, address, uint16(len(values)/2), values)
	if err != nil {
		c.handleConnectionError(err)
		return err
	}
	return nil
}

// WriteSingleCoil writes a single coil to the MODBUS server
func (c *Client) WriteSingleCoil(ctx context.Context, address uint16, value uint16) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("modbus client not connected")
	}

	_, err := c.client.WriteSingleCoil(ctx, address, value)
	if err != nil {
		c.handleConnectionError(err)
		return err
	}
	return nil
}

// WriteMultipleCoils writes multiple coils to the MODBUS server
func (c *Client) WriteMultipleCoils(ctx context.Context, address, quantity uint16, values []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("modbus client not connected")
	}

	_, err := c.client.WriteMultipleCoils(ctx, address, quantity, values)
	if err != nil {
		c.handleConnectionError(err)
		return err
	}
	return nil
}

// withSlaveID executes a function with a temporary slave ID, then restores the original
func (c *Client) withSlaveID(slaveID byte, fn func() error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Change slave ID
	originalSlaveID := c.handler.SlaveID
	c.handler.SlaveID = slaveID

	// Execute function
	err := fn()

	// Restore slave ID
	c.handler.SlaveID = originalSlaveID
	return err
}

// ReadHoldingRegistersWithSlaveID reads holding registers with a specific slave ID
func (c *Client) ReadHoldingRegistersWithSlaveID(ctx context.Context, slaveID byte, address, quantity uint16) ([]byte, error) {
	var result []byte

	err := c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		data, err := c.client.ReadHoldingRegisters(ctx, address, quantity)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
			return err
		}
		result = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReadInputRegistersWithSlaveID reads input registers with a specific slave ID
func (c *Client) ReadInputRegistersWithSlaveID(ctx context.Context, slaveID byte, address, quantity uint16) ([]byte, error) {
	var result []byte

	err := c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		data, err := c.client.ReadInputRegisters(ctx, address, quantity)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
			return err
		}
		result = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReadDiscreteInputsWithSlaveID reads discrete inputs with a specific slave ID
func (c *Client) ReadDiscreteInputsWithSlaveID(ctx context.Context, slaveID byte, address, quantity uint16) ([]byte, error) {
	var result []byte

	err := c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		data, err := c.client.ReadDiscreteInputs(ctx, address, quantity)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
			return err
		}
		result = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReadCoilsWithSlaveID reads coils with a specific slave ID
func (c *Client) ReadCoilsWithSlaveID(ctx context.Context, slaveID byte, address, quantity uint16) ([]byte, error) {
	var result []byte

	err := c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		data, err := c.client.ReadCoils(ctx, address, quantity)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
			return err
		}
		result = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// WriteSingleRegisterWithSlaveID writes a single register with a specific slave ID
func (c *Client) WriteSingleRegisterWithSlaveID(ctx context.Context, slaveID byte, address, value uint16) error {
	return c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		_, err := c.client.WriteSingleRegister(ctx, address, value)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
		}
		return err
	})
}

// WriteMultipleRegistersWithSlaveID writes multiple registers with a specific slave ID
func (c *Client) WriteMultipleRegistersWithSlaveID(ctx context.Context, slaveID byte, address uint16, values []byte) error {
	return c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		if len(values)%2 != 0 {
			return fmt.Errorf("values must be even number of bytes, got %d", len(values))
		}

		_, err := c.client.WriteMultipleRegisters(ctx, address, uint16(len(values)/2), values)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
		}
		return err
	})
}

// WriteSingleCoilWithSlaveID writes a single coil with a specific slave ID
func (c *Client) WriteSingleCoilWithSlaveID(ctx context.Context, slaveID byte, address, value uint16) error {
	return c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		_, err := c.client.WriteSingleCoil(ctx, address, value)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
		}
		return err
	})
}

// WriteMultipleCoilsWithSlaveID writes multiple coils with a specific slave ID
func (c *Client) WriteMultipleCoilsWithSlaveID(ctx context.Context, slaveID byte, address, quantity uint16, values []byte) error {
	return c.withSlaveID(slaveID, func() error {
		if !c.isConnected {
			return fmt.Errorf("modbus client not connected")
		}

		_, err := c.client.WriteMultipleCoils(ctx, address, quantity, values)
		if err != nil {
			if !c.isModbusProtocolError(err) {
				c.isConnected = false
			}
		}
		return err
	})
}

// handleConnectionError checks if the error indicates a connection loss and updates the flag
func (c *Client) handleConnectionError(err error) {
	if err != nil && !c.isModbusProtocolError(err) {
		go c.markDisconnected()
	}
}

// isModbusProtocolError determines if an error is a valid Modbus protocol error
func (c *Client) isModbusProtocolError(err error) bool {
	var modbusErr *modbus.Error
	return errors.As(err, &modbusErr)
}

// markDisconnected safely marks the client as disconnected
func (c *Client) markDisconnected() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isConnected = false
}
