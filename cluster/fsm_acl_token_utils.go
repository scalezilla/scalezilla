package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/Lord-Y/rafty"
)

// aclTokenEncodeCommand permits to transform command receive from clients to binary language machine
func aclTokenEncodeCommand(cmd aclTokenCommand, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(cmd.Kind)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint64(len(cmd.AccessorID))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(cmd.AccessorID)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint64(len(cmd.Token))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(cmd.Token)); err != nil {
		return err
	}
	// the following is ONLY used for boolean
	if err := binary.Write(w, binary.LittleEndian, cmd.InitialToken); err != nil {
		return err
	}
	return nil
}

// aclTokenDecodeCommand permits to transform back command from binary language machine to clients
func aclTokenDecodeCommand(data []byte) (aclTokenCommand, error) {
	var cmd aclTokenCommand
	buffer := bytes.NewBuffer(data)

	var kind uint32
	if err := binary.Read(buffer, binary.LittleEndian, &kind); err != nil {
		return cmd, err
	}
	cmd.Kind = commandKind(kind)

	var accessorIDLen uint64
	if err := binary.Read(buffer, binary.LittleEndian, &accessorIDLen); err != nil {
		return cmd, err
	}
	accessorID := make([]byte, accessorIDLen)
	if _, err := buffer.Read(accessorID); err != nil {
		return cmd, err
	}
	cmd.AccessorID = string(accessorID)

	var tokenLen uint64
	if err := binary.Read(buffer, binary.LittleEndian, &tokenLen); err != nil {
		return cmd, err
	}
	token := make([]byte, tokenLen)
	if _, err := buffer.Read(token); err != nil {
		return cmd, err
	}
	cmd.Token = string(token)

	var initialtokenLen bool
	if err := binary.Read(buffer, binary.LittleEndian, &initialtokenLen); err != nil {
		return cmd, err
	}
	cmd.InitialToken = initialtokenLen

	return cmd, nil
}

// aclTokenApplyCommand will apply fsm to the acl token store
func (f *fsmState) aclTokenApplyCommand(log *rafty.LogEntry) ([]byte, error) {
	cmd, _ := aclTokenDecodeCommand(log.Command)

	if rafty.LogKind(log.LogType) == rafty.LogCommandReadLeaderLease || rafty.LogKind(log.LogType) == rafty.LogCommandLinearizableRead {
		return f.memoryStore.aclTokenEncoded(cmd)
	}

	switch cmd.Kind {
	case aclTokenCommandSet:
		return nil, f.memoryStore.aclTokenSet(log, cmd)

	case aclTokenCommandGet:
		value, err := f.memoryStore.aclTokenGet([]byte(cmd.AccessorID))
		if err != nil {
			return nil, err
		}
		return value, nil

	case aclTokenCommandDelete:
		f.memoryStore.aclTokenDelete([]byte(cmd.AccessorID))
	}

	return nil, nil
}

// aclTokenUnmarshal is an helper to unmarshal bytes into acl token struct
func aclTokenUnmarshal(data []byte) (AclToken, error) {
	var token AclToken
	if err := json.Unmarshal(data, &token); err != nil {
		return token, err
	}
	return token, nil
}
