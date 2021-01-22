package ffiwrapper

import (
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
)

type Config struct {
	SealProofType abi.RegisteredSealProof

	_ struct{} // guard against nameless init
}

func sizeFromConfig(cfg Config) (abi.SectorSize, error) {
	return cfg.SealProofType.SectorSize()
}

func SealProofTypeFromSectorSize(ssize abi.SectorSize) (abi.RegisteredSealProof, error) {
	switch ssize {
	case 2 << 10:
		return abi.RegisteredSealProof_StackedDrg2KiBV1, nil
	case 8 << 20:
		return abi.RegisteredSealProof_StackedDrg8MiBV1, nil
	case 512 << 20:
		return abi.RegisteredSealProof_StackedDrg512MiBV1, nil
	case 32 << 30:
		return abi.RegisteredSealProof_StackedDrg32GiBV1, nil
	case 64 << 30:
		return abi.RegisteredSealProof_StackedDrg64GiBV1, nil
	case 4 << 30:
		return abi.RegisteredSealProof_StackedDrg4GiBV1, nil
	case 16 << 30:
		return abi.RegisteredSealProof_StackedDrg16GiBV1, nil
	default:
		return 0, xerrors.Errorf("unsupported sector size for miner: %v", ssize)
	}
}

func SealProofTypeFromSectorSizeForAddSector(ssize abi.SectorSize, nv network.Version) (abi.RegisteredSealProof, error) {
	switch {
	case nv < network.Version4:
		switch ssize {
		case 2 << 10:
			return abi.RegisteredSealProof_StackedDrg2KiBV1, nil
		case 8 << 20:
			return abi.RegisteredSealProof_StackedDrg8MiBV1, nil
		case 512 << 20:
			return abi.RegisteredSealProof_StackedDrg512MiBV1, nil
		case 32 << 30:
			return abi.RegisteredSealProof_StackedDrg32GiBV1, nil
		case 64 << 30:
			return abi.RegisteredSealProof_StackedDrg64GiBV1, nil
		case 4 << 30:
			return abi.RegisteredSealProof_StackedDrg4GiBV1, nil
		default:
			return 0, xerrors.Errorf("unsupported sector size for miner: %v", ssize)
		}
	case nv >= network.Version4:
		switch ssize {
		case 2 << 10:
			return abi.RegisteredSealProof_StackedDrg2KiBV1, nil
		case 8 << 20:
			return abi.RegisteredSealProof_StackedDrg8MiBV1, nil
		case 512 << 20:
			return abi.RegisteredSealProof_StackedDrg512MiBV1, nil
		case 32 << 30:
			return abi.RegisteredSealProof_StackedDrg32GiBV1, nil
		case 64 << 30:
			return abi.RegisteredSealProof_StackedDrg64GiBV1, nil
		case 4 << 30:
			return abi.RegisteredSealProof_StackedDrg4GiBV1, nil
		case 16 << 30:
			return abi.RegisteredSealProof_StackedDrg16GiBV1, nil
		default:
			return 0, xerrors.Errorf("unsupported sector size for miner: %v", ssize)
		}
	}
	return 0, xerrors.Errorf("unsupported network version")
}