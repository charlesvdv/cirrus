package sqlite

import (
	"strconv"

	"crawshaw.io/sqlite"
	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database"
	"github.com/rs/zerolog/log"
)

// FilesystemRepository contains the sqlite persistence layer for filesystem related data
type FilesystemRepository struct {
}

// CreateDirectory impl
func (r FilesystemRepository) CreateDirectory(_tx database.Tx, directory *cirrus.Directory) error {
	tx := getTx(_tx)

	stmt := tx.Prep("INSERT INTO directories (parent_id, owner_id, name, created_at) VALUES ($parent_id, $owner_id, $name, $created_at);")
	if directory.Parent == nil {
		stmt.SetNull("$parent_id")
	} else {
		stmt.SetInt64("$parent_id", formatObjectID(*directory.Parent))
	}
	stmt.SetInt64("$owner_id", formatUserID(directory.Owner))
	stmt.SetText("$name", directory.Name)
	stmt.SetText("$created_at", formatTime(directory.CreatedAt))
	_, err := stmt.Step()
	defer stmt.Reset()
	if err != nil {
		log.Err(err).Msg("failed to create directory")
		return formatError(err)
	}

	directory.ID = cirrus.ObjectID(strconv.FormatInt(tx.LastInsertRowID(), 10))

	return nil
}

// ListDirectoryContent impl
func (r FilesystemRepository) ListDirectoryContent(_tx database.Tx, ownerID cirrus.UserID, parent *cirrus.ObjectID) ([]cirrus.FilesystemObject, error) {
	tx := getTx(_tx)

	var objects []cirrus.FilesystemObject
	var directory cirrus.Directory

	stmt := tx.Prep("SELECT id, parent_id, owner_id, name, created_at FROM directories WHERE parent_id IS $parent_id AND owner_id = $owner_id;")
	if parent == nil {
		stmt.SetNull("$parent_id")
	} else {
		stmt.SetInt64("$parent_id", formatObjectID(*parent))
	}
	stmt.SetInt64("$owner_id", formatUserID(ownerID))
	for {
		hasRow, err := stmt.Step()
		defer stmt.Reset()
		if err != nil {
			log.Err(err).Msg("failed to list directory content")
			return nil, formatError(err)
		}
		if !hasRow {
			break
		}
		directory.ID = getObjectID(stmt.GetInt64("id"))
		if stmt.GetType("parent_id") == sqlite.SQLITE_NULL {
			directory.Parent = nil
		} else {
			parentID := getObjectID(stmt.GetInt64("parent_id"))
			directory.Parent = &parentID
		}
		directory.Owner = getUserID(stmt.GetInt64("owner_id"))
		directory.Name = stmt.GetText("name")
		directory.CreatedAt = getTime(stmt.GetText("created_at"))

		objects = append(objects, directory)
	}

	return objects, nil
}

// ResolvePath impl
func (r FilesystemRepository) ResolvePath(_tx database.Tx, ownerID cirrus.UserID, path cirrus.Path) ([]cirrus.FilesystemObject, error) {
	tx := getTx(_tx)

	// TODO: there is probably a better sql way but for now, it's good enough™
	var parentID int64
	var fsObjects []cirrus.FilesystemObject

	dirStmt := tx.Prep("SELECT id, parent_id, owner_id, name, created_at FROM directories WHERE name = $name AND parent_id IS $parent_id AND owner_id = $owner_id")
	dirStmt.SetInt64("$owner_id", formatUserID(ownerID))
	lastElemIsMaybeAFile := false
	for pathIndex, elem := range path.Elements() {
		dirStmt.SetText("$name", elem)
		if pathIndex == 0 {
			dirStmt.SetNull("$parent_id")
		} else {
			dirStmt.SetInt64("$parent_id", parentID)
		}

		hasRow, err := dirStmt.Step()
		defer dirStmt.Reset()
		if err != nil {
			log.Err(err).Msg("failed to fetch directories info")
			return nil, formatError(err)
		}

		if !hasRow {
			if pathIndex == path.ElementCount()-1 {
				// We have the last element not resolved yet... It may be a file
				lastElemIsMaybeAFile = true
			}
			break
		}

		var directory cirrus.Directory
		directory.ID = getObjectID(dirStmt.GetInt64("id"))
		if dirStmt.GetType("parent_id") == sqlite.SQLITE_NULL {
			directory.Parent = nil
		} else {
			parent := getObjectID(dirStmt.GetInt64("parent_id"))
			directory.Parent = &parent
		}
		directory.Owner = getUserID(dirStmt.GetInt64("owner_id"))
		directory.Name = dirStmt.GetText("name")
		directory.CreatedAt = getTime(dirStmt.GetText("created_at"))

		fsObjects = append(fsObjects, directory)
		parentID = formatObjectID(directory.ID)
	}

	if lastElemIsMaybeAFile {
		// TODO
	}

	return fsObjects, nil
}
