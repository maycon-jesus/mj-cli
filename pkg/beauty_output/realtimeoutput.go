package beautyoutput

import "fmt"

func (sb *StrBuilder) IsRealtimeOutput() bool {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.RealtimeOutput
}

func (sb *StrBuilder) SetRealtimeOutput(realtime bool) *StrBuilder {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.RealtimeOutput = realtime
	return sb
}

func (sb *StrBuilder) printRealtimeOutputUnsafe() *StrBuilder {
	if !sb.RealtimeOutput {
		return sb
	}
	fmt.Print(sb.content.String())
	sb.content.Reset()
	return sb
}
