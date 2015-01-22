package booking

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// Locator for a guest record
type guestId uint32

func (id *guestId) Scan(src interface{}) error {
	n, ok := src.(int64)
	if !ok {
		b, _ := src.([]byte)
		err := errors.New(
			fmt.Sprintf("can't scan guestId from db: %s", string(b)),
		)
		glog.Error(err)
		return err
	}
	*id = guestId(n)
	return nil
}

func (id guestId) Value() (driver.Value, error) {
	return driver.Value(int64(id)), nil
}

func (id guestId) String() string {
	return fmt.Sprintf("guestId:%d", id)
}

// A guest record
type guest struct {
	Email       email
	Id          guestId
	Name        name
	PhoneNumber phoneNumber
}

// List of registered guests
type Guestbook struct {
	DB *sql.DB `inject:""`

	lookup       *sql.Stmt
	register     *sql.Stmt
	tableCreated bool
}

var guestNotFound = errors.New("guest not found")

func (g *Guestbook) createTableOnce() {
	if g.tableCreated {
		return
	}

	_, err := g.DB.Exec(`
    create table Guestbook (
      Email text unique not null,
      Id integer primary key autoincrement not null,
      Name text not null,
      PhoneNumber text not null
    )
  `)
	if err == nil {
		g.tableCreated = true
		glog.Info("Guestbook table created")
	} else {
		glog.Warning(err)
	}
}

func (g *Guestbook) Lookup(byId guestId) (guest, error) {
	g.createTableOnce()

	if g.lookup == nil {
		var err error
		g.lookup, err = g.DB.Prepare(
			`select Email, Id, Name, PhoneNumber from Guestbook where Id = $1`,
		)
		if err != nil {
			return guest{}, err
		}
	}

	var row_email string
	var row_id int
	var row_name string
	var row_phone_number string

	row := g.lookup.QueryRow(byId)
	err := row.Scan(
		&row_email,
		&row_id,
		&row_name,
		&row_phone_number,
	)

	if err == sql.ErrNoRows {
		return guest{}, guestNotFound
	}

	if err != nil {
		return guest{}, err
	}

	var email email
	if !email.Set(row_email) {
		return guest{}, err
	}

	var name name
	if !name.Set(row_name) {
		return guest{}, nil
	}

	return guest{
		Email:       email,
		Id:          guestId(row_id),
		Name:        name,
		PhoneNumber: phoneNumber(row_phone_number),
	}, nil
}

// If email already registered, return existing record
// If someone wants to sign someone else up and pay, fine
func (g *Guestbook) Register(
	name *name,
	email *email,
	phone phoneNumber,
) (guestId, error) {
	g.createTableOnce()

	if g.register == nil {
		var err error
		g.register, err = g.DB.Prepare(`
      insert into Guestbook (Email, Name, PhoneNumber)
      values ($1, $2, $3)
    `)
		if err != nil {
			return 0, err
		}
	}

	result, err := g.register.Exec(email.s, name, phone)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return guestId(lastId), nil
}

// A guest's name - must be 2 parts, not empty
type name struct {
	s string
}

func (n *name) Value() (driver.Value, error) {
	return driver.Value(n.s), nil
}

func (n *name) Equal(s string) bool {
	return n.s == s
}

func (n *name) Set(s string) bool {
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return false
	}

	n.s = strings.Join(fields, " ")
	return true
}

// An electronic mail address
type email struct {
	s string
}

func (e *email) Value() (driver.Value, error) {
	return driver.Value(e.s), nil
}

func (e *email) Equal(s string) bool {
	return e.s == s
}

func (e *email) Set(s string) bool {
	s = strings.Trim(s, " ")

	if len(s) < 3 {
		return false
	}

	if !strings.Contains(s, "@") {
		return false
	}

	e.s = s
	return true
}

// A telephone number
type phoneNumber string

func (p phoneNumber) Value() (driver.Value, error) {
	return driver.Value(string(p)), nil
}
