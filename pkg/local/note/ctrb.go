package note

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/notebox/nb-crdt-go/block"
	"github.com/notebox/nb-crdt-go/common"
)

func Prepare(db *sql.DB, path string) error {
	ok, err := isInstalled(db)
	if err != nil {
		return err
	}

	if !ok {
		err = addDemoNote(path)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS nb_note_ctrbs (
			note_id TEXT NOT NULL,
			block_id TEXT NOT NULL,
			block_nonce INTEGER NOT NULL,
			text_nonce INTEGER NOT NULL,
			replica_id INTEGER NOT NULL,
			timestamp INTEGER NOT NULL,
			ops BLOB NOT NULL,
			UNIQUE (note_id, block_id, replica_id, block_nonce, text_nonce)
		);
		CREATE INDEX IF NOT EXISTS idx_nb_note_ctrbs_ids ON nb_note_ctrbs (note_id, block_id);
		CREATE INDEX IF NOT EXISTS idx_nb_note_ctrbs_block_nonce ON nb_note_ctrbs (block_nonce);
		CREATE TABLE IF NOT EXISTS nb_note_blocks (
			note_id TEXT NOT NULL,
			block_id TEXT NOT NULL,
			data BLOB NOT NULL,
			UNIQUE (note_id, block_id)
		);
		CREATE INDEX IF NOT EXISTS idx_nb_note_blocks_ids ON nb_note_blocks (note_id, block_id);
	`)
	return err
}

func isInstalled(db *sql.DB) (bool, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='nb_note_ctrbs'")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func InsertCTRBs(db *sql.DB, noteID *uuid.UUID, contributions []*block.Contribution) ([]*uuid.UUID, error) {
	var blockIDs []*uuid.UUID
	q := "INSERT INTO nb_note_ctrbs (note_id, block_id, block_nonce, text_nonce, replica_id, timestamp, ops) VALUES"
	v := []any{}

	for _, row := range contributions {
		blockIDs = append(blockIDs, &row.BlockID)
		q += "(?, ?, ?, ?, ?, ?, ?),"
		data, err := json.Marshal(row.Operations)
		if err != nil {
			return nil, err
		}
		v = append(v, noteID, row.BlockID, row.Nonce[0], row.Nonce[1], row.Stamp.ReplicaID, row.Stamp.Timestamp, data)
	}

	stmt, err := db.Prepare(q[:len(q)-1])
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(v...)

	return blockIDs, err
}

func InsertBlock(db *sql.DB, noteID *uuid.UUID, block *block.Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT OR REPLACE INTO nb_note_blocks (note_id, block_id, data) VALUES (?, ?, ?)", noteID, block.BlockID, data)
	return err
}

func SelectBlockData(db *sql.DB, noteID *uuid.UUID, blockID *uuid.UUID) ([]byte, error) {
	rows, err := db.Query("SELECT data FROM nb_note_blocks WHERE note_id = ? AND block_id LIMIT 1", noteID, blockID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func SelectBlocks(db *sql.DB, noteID *uuid.UUID) ([]*block.Block, error) {
	var blocks []*block.Block

	rows, err := db.Query("SELECT data FROM nb_note_blocks WHERE note_id = ?", noteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return blocks, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		var b block.Block
		err = json.Unmarshal(data, &b)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &b)
	}

	return blocks, nil
}

func SelectAllAfter(db *sql.DB, replicaID uint32, noteID *uuid.UUID, blockID *uuid.UUID, blockNonce common.Nonce, textNonce common.Nonce) ([]*block.Contribution, error) {
	rows, err := db.Query("SELECT block_id, block_nonce, text_nonce, replica_id, timestamp, ops FROM nb_note_ctrbs WHERE replica_id = ? AND note_id = ? AND block_id = ? AND (block_nonce > ? OR (block_nonce = ? AND text_nonce > ?)) ORDER BY block_nonce ASC", replicaID, noteID, blockID, blockNonce, blockNonce, textNonce)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ctrbs []*block.Contribution
	for rows.Next() {
		var ctrb block.Contribution
		ctrb.Nonce = make([]common.Nonce, 2)
		var data []byte
		err := rows.Scan(&ctrb.BlockID, &ctrb.Nonce[0], &ctrb.Nonce[1], &ctrb.Stamp.ReplicaID, &ctrb.Stamp.Timestamp, &data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &ctrb.Operations)
		if err != nil {
			return nil, err
		}
		ctrbs = append(ctrbs, &ctrb)
	}

	if err == sql.ErrNoRows {
		err = nil
	}

	return ctrbs, err
}
