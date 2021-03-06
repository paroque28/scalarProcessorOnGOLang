package cpu

import "fmt"

type WriteBack struct {
	InMemControlSignals MemControlSignals `json:"in_mem_control_signals"`
	InALUResult         uint64            `json:"in_alu_result"`
}

func (wb *WriteBack) Run(done chan string, registers []uint64) {
	if wb.InMemControlSignals.RegisterWriteEnable {
		registers[wb.InMemControlSignals.WriteAddress] = wb.InALUResult
		if DEBUG > 3 {
			fmt.Printf("[WriteBack] V[%x] = %x\n", wb.InMemControlSignals.WriteAddress, wb.InALUResult)
		}
	}
	done <- "wb"
}

func (wb *WriteBack) UpdateInRegisters(inControlSignals MemControlSignals, ALUResult uint64) {
	wb.InALUResult = ALUResult
	wb.InMemControlSignals = inControlSignals
}
