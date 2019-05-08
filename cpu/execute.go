package cpu

import "fmt"

type Execute struct {
	InDecodeControlSignals DecodeControlSignals  `json:"in_control_signals"`
	InRd1                  uint64                `json:"in_rd_1"`
	InRd2                  uint64                `json:"in_rd_2"`
	InImmediate            int64                 `json:"in_immediate"`
	ALUResult              uint64                `json:"alu_result"`
	OutControlSignals      ExecuteControlSignals `json:"control_signals"`
}

type ExecuteControlSignals struct {
	MemWriteAddress     uint64 `json:"mem_write_address"`
	WriteAddress        byte   `json:"write_address"`
	MemWriteEnable      bool   `json:"mem_write_enable"`
	MemToReg            bool   `json:"mem_to_reg"`
	RegisterWriteEnable bool   `json:"register_write_enable"`
}

const (
	ALU_NOP       = iota
	ALU_BUFFER    = iota
	ALU_ADD       = iota
	ALU_XOR       = iota
	ALU_TOTAL_OPS = iota
)

func (exec *Execute) Run(done chan string) {
	if exec.InDecodeControlSignals.MemWriteAddress != 0 {
		fmt.Println("mem")
	}
	if exec.InDecodeControlSignals.ALUSrcReg {
		exec.ALUResult = ALU(exec.InDecodeControlSignals.ALUControl, exec.InRd1, exec.InRd2)
	} else {
		exec.ALUResult = ALU(exec.InDecodeControlSignals.ALUControl, exec.InRd1, uint64(exec.InImmediate))
	}
	if exec.InDecodeControlSignals.ALUControl != 0 {
		fmt.Printf("[Exec] Rd1: %x Op:%x Rd2: %x, Imm:%x, ALUResult = %x\n", exec.InRd1, exec.InDecodeControlSignals.ALUControl, exec.InRd2, exec.InImmediate, exec.ALUResult)
	}
	exec.setControlSignals(exec.InDecodeControlSignals.MemWriteEnable,
		exec.InDecodeControlSignals.MemToReg,
		exec.InDecodeControlSignals.RegisterWriteEnable,
		exec.InDecodeControlSignals.WriteAddress)
	done <- "execute"
}
func ALU(aluOp byte, a uint64, b uint64) (result uint64) {
	switch aluOp {
	case ALU_BUFFER:
		result = a
	case ALU_NOP:
		result = 0
	case ALU_ADD:
		result = uint64(int(a) + int(b))
	default:
		panic("ALU operation not implemented")
	}
	return
}
func (exec *Execute) UpdateInRegisters(inControlSignals DecodeControlSignals, rd1 uint64, rd2 uint64, immediate int64) {
	exec.InDecodeControlSignals = inControlSignals
	exec.InRd1 = rd1
	exec.InRd2 = rd2
	exec.InImmediate = immediate
}

func (exec *Execute) setControlSignals(memWriteEnable bool,
	memToReg bool,
	registerWriteEnable bool,
	writeAddress byte) {

	exec.OutControlSignals.MemToReg = memToReg
	exec.OutControlSignals.MemWriteEnable = memWriteEnable
	exec.OutControlSignals.RegisterWriteEnable = registerWriteEnable
	exec.OutControlSignals.WriteAddress = writeAddress
	exec.OutControlSignals.MemWriteAddress = exec.InDecodeControlSignals.MemWriteAddress
}
