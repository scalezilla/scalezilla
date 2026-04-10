package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"maps"
	"slices"
	"time"

	"github.com/Lord-Y/rafty"
)

// deploymentEncodeCommand permits to transform command receive from clients to binary language machine
func deploymentEncodeCommand(cmd deploymentState, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(cmd.Kind)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint64(len(cmd.Name))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(cmd.Name)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, cmd.NewRollingVersion); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, cmd.CurrentUsedVersion); err != nil {
		return err
	}

	keys := slices.Sorted(maps.Keys(cmd.Content))
	// count is used to know the length of the slice to be used later for decoding
	if err := binary.Write(w, binary.LittleEndian, uint64(len(keys))); err != nil {
		return err
	}
	for _, k := range keys {
		if err := binary.Write(w, binary.LittleEndian, k); err != nil {
			return err
		}
		v := cmd.Content[k]
		if err := binary.Write(w, binary.LittleEndian, v.IsStable); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint64(len(v.RawContent))); err != nil {
			return err
		}
		if _, err := w.Write([]byte(v.RawContent)); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, v.Version); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, v.CreatedAt.UnixNano()); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint64(len(v.ReplicaSetID))); err != nil {
			return err
		}
		if _, err := w.Write([]byte(v.ReplicaSetID)); err != nil {
			return err
		}
	}
	// the following is ONLY used for boolean
	if err := binary.Write(w, binary.LittleEndian, cmd.MustBeStarted); err != nil {
		return err
	}
	return nil
}

// deploymentDecodeCommand permits to transform back command from binary language machine to clients
func deploymentDecodeCommand(data []byte) (deploymentState, error) {
	var cmd deploymentState
	buffer := bytes.NewBuffer(data)

	var kind uint32
	if err := binary.Read(buffer, binary.LittleEndian, &kind); err != nil {
		return cmd, err
	}
	cmd.Kind = commandKind(kind)

	var nameLen uint64
	if err := binary.Read(buffer, binary.LittleEndian, &nameLen); err != nil {
		return cmd, err
	}
	name := make([]byte, nameLen)
	if _, err := io.ReadFull(buffer, name); err != nil {
		return cmd, err
	}
	cmd.Name = string(name)

	var newRollingVersion int64
	if err := binary.Read(buffer, binary.LittleEndian, &newRollingVersion); err != nil {
		return cmd, err
	}
	cmd.NewRollingVersion = newRollingVersion

	var currentUsedVersion uint64
	if err := binary.Read(buffer, binary.LittleEndian, &currentUsedVersion); err != nil {
		return cmd, err
	}
	cmd.CurrentUsedVersion = currentUsedVersion

	// content part
	var count uint64
	if err := binary.Read(buffer, binary.LittleEndian, &count); err != nil {
		return cmd, err
	}

	cmd.Content = make(map[uint64]deploymentContent, count)
	for range count {
		var content deploymentContent
		var key uint64
		if err := binary.Read(buffer, binary.LittleEndian, &key); err != nil {
			return cmd, err
		}

		var isStable bool
		if err := binary.Read(buffer, binary.LittleEndian, &isStable); err != nil {
			return cmd, err
		}
		content.IsStable = isStable

		var rawContentLen uint64
		if err := binary.Read(buffer, binary.LittleEndian, &rawContentLen); err != nil {
			return cmd, err
		}
		rawContent := make([]byte, rawContentLen)
		if _, err := io.ReadFull(buffer, rawContent); err != nil {
			return cmd, err
		}
		content.RawContent = string(rawContent)

		var version uint64
		if err := binary.Read(buffer, binary.LittleEndian, &version); err != nil {
			return cmd, err
		}
		content.Version = version

		var createdAt int64
		if err := binary.Read(buffer, binary.LittleEndian, &createdAt); err != nil {
			return cmd, err
		}
		content.CreatedAt = time.Unix(0, createdAt)

		var replicaSetIDLen uint64
		if err := binary.Read(buffer, binary.LittleEndian, &replicaSetIDLen); err != nil {
			return cmd, err
		}
		replicaSetID := make([]byte, replicaSetIDLen)
		if _, err := io.ReadFull(buffer, replicaSetID); err != nil {
			return cmd, err
		}
		content.ReplicaSetID = string(replicaSetID)

		cmd.Content[key] = content
	}

	var mustBeStarted bool
	if err := binary.Read(buffer, binary.LittleEndian, &mustBeStarted); err != nil {
		return cmd, err
	}
	cmd.MustBeStarted = mustBeStarted

	return cmd, nil
}

// deploymentApplyCommand will apply fsm to the deployment store
func (f *fsmState) deploymentApplyCommand(log *rafty.LogEntry) ([]byte, error) {
	cmd, _ := deploymentDecodeCommand(log.Command)

	if rafty.LogKind(log.LogType) == rafty.LogCommandReadLeaderLease || rafty.LogKind(log.LogType) == rafty.LogCommandLinearizableRead {
		return f.memoryStore.deploymentEncoded(cmd)
	}

	switch cmd.Kind {
	case deploymentCommandSet:
		return nil, f.memoryStore.deploymentSet(log, cmd)

	case deploymentCommandGet:
		value, err := f.memoryStore.deploymentGet([]byte(cmd.Name))
		if err != nil {
			return nil, err
		}
		return value, nil

	case deploymentCommandDelete:
		f.memoryStore.deploymentDelete([]byte(cmd.Name))
	}

	return nil, nil
}

// deploymentUnmarshal is an helper to unmarshal bytes into deploymentState struct
func deploymentUnmarshal(data []byte) (deploymentState, error) {
	var d deploymentState
	if err := json.Unmarshal(data, &d); err != nil {
		return d, err
	}
	return d, nil
}
