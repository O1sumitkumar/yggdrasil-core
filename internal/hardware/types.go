package hardware

import (
	"context"
	"runtime"

	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

func detectCPU(ctx context.Context) (contracts.CPUInfo, error) {
	info, err := platformCPU(ctx)
	if err != nil {
		return contracts.CPUInfo{
			Model:   "Unknown CPU",
			Cores:   runtime.NumCPU(),
			Threads: runtime.NumCPU(),
		}, err
	}
	if info.Cores == 0 {
		info.Cores = runtime.NumCPU()
	}
	if info.Threads == 0 {
		info.Threads = runtime.NumCPU()
	}
	return info, nil
}

func detectMemory(ctx context.Context) (contracts.MemoryInfo, error) {
	return platformMemory(ctx)
}

func detectDisk(ctx context.Context, path string) (contracts.DiskInfo, error) {
	return platformDisk(ctx, path)
}

func detectAccelerators(ctx context.Context) ([]contracts.Accelerator, error) {
	return platformAccelerators(ctx)
}
