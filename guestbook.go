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
			panic(err)
		}
	}

	var id guestId
	var name name
	var email email
	var phone phoneNumber

	row := g.lookup.QueryRow(byId)
	err := row.Scan(
		&email,
		&id,
		&name,
		&phone,
	)

	if err == sql.ErrNoRows {
		return guest{}, guestNotFound
	}

	if err != nil {
		glog.Error(err)
		return guest{}, err
	}

	return guest{
		Email:       email,
		Id:          id,
		Name:        name,
		PhoneNumber: phone,
	}, nil
}

// If email already registered, return existing record
// If someone wants to sign someone else up and pay, fine
func (g *Guestbook) Register(
	name name,
	email email,
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
			panic(err)
		}
	}

	result, err := g.register.Exec(email, name, phone)
	if err != nil {
		glog.Error(err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
		return 0, err
	}

	return guestId(lastId), nil
}

// A guest's name - must be 2 parts, not empty
type name struct {
	s string
}

var invalidName = errors.New("invalid guest name")

func newName(s string) (name, error) {
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return name{}, invalidName
	}

	return name{s: strings.Join(fields, " ")}, nil
}

func (n *name) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan name from %#v", src),
		)
		glog.Error(err)
		return err
	}

	name, err := newName(string(b))
	if err != nil {
		return err
	}
	n.s = name.s
	return nil
}

func (n name) Value() (driver.Value, error) {
	return driver.Value(n.s), nil
}

func (n name) Equal(s string) bool {
	return n.s == s
}

// An electronic mail address
type email struct {
	s string
}

var invalidEmail = errors.New("invalid email")

func newEmail(s string) (email, error) {
	s = strings.Trim(s, " ")

	if len(s) < 3 {
		return email{}, invalidEmail
	}

	if !strings.Contains(s, "@") {
		return email{}, invalidEmail
	}

	return email{s: s}, nil
}

func (e *email) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan email from %#v", src),
		)
		glog.Error(err)
		return err
	}

	email, err := newEmail(string(b))
	if err != nil {
		return err
	}
	e.s = email.s
	return nil
}

func (e email) Value() (driver.Value, error) {
	return driver.Value(e.s), nil
}

func (e email) Equal(s string) bool {
	return e.s == s
}

// A telephone number
type phoneNumber struct {
	s string
}

// TODO validate, maybe libphonenumber?
func newPhoneNumber(s string) (phoneNumber, error) {
	return phoneNumber{s: s}, nil
}

func (p *phoneNumber) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		err := errors.New(
			fmt.Sprintf("can't scan phoneNumber from %#v", src),
		)
		glog.Error(err)
		return err
	}

	phone, err := newPhoneNumber(string(b))
	if err != nil {
		return err
	}
	p.s = phone.s
	return nil
}

func (p phoneNumber) Value() (driver.Value, error) {
	return driver.Value(p.s), nil
}
