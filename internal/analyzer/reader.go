package analyzer

import "fmt"

// readBaseData reads the base data from the energy analyzer
func (s *Service) readBaseData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, BaseDataStartAddr, BaseDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	baseData := parseBaseData(data)

	s.mutex.Lock()
	s.lastData = baseData
	s.mutex.Unlock()

	return nil
}
