package booking

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

// Guestbook is list of registered guests
type Guestbook struct {
	DB *sql.DB `inject:""`
}

const GuestbookSchema = `
  CREATE TABLE Guestbook (
    Email TEXT UNIQUE NOT NULL,
    Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Name TEXT NOT NULL,
    PhoneNumber TEXT NOT NULL
  )
`

var guestNotFound = errors.New("guest not found")

func (g *Guestbook) Lookup(byId guestId) (guest, error) {
	var guest guest
	err := g.withTx(func(tx *GuestbookTx) error {
		var err error
		guest, err = tx.Lookup(byId)
		return err
	})
	return guest, err
}

// Emails are unique
func (g *Guestbook) Register(
	name name,
	email email,
	phone phoneNumber,
) (guestId, error) {
	var id guestId
	err := g.withTx(func(tx *GuestbookTx) error {
		var err error
		id, err = tx.Register(name, email, phone)
		return err
	})
	return id, err
}

// Wraps fn to provide a CalendarTx that calls Rollback on err, Commit on ok
// Theoretically you could call a series of funcs for batch
func (g *Guestbook) withTx(fn func(*GuestbookTx) error) error {
	dbTx, err := g.DB.Begin()
	if err != nil {
		return err
	}
	tx := &GuestbookTx{dbTx}
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type GuestbookTx struct {
	*sql.Tx
}

func (tx *GuestbookTx) Lookup(byId guestId) (guest, error) {
	stmt, err := tx.Prepare(`
		select Email, Id, Name, PhoneNumber 
		from Guestbook 
		where Id = $1`,
	)
	if err != nil {
		panic(err)
	}

	var id guestId
	var name name
	var email email
	var phone phoneNumber

	row := stmt.QueryRow(byId)
	err = row.Scan(
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

func (tx *GuestbookTx) Register(
	name name,
	email email,
	phone phoneNumber,
) (guestId, error) {
	stmt, err := tx.Prepare(`
      insert into Guestbook (Email, Name, PhoneNumber)
      values ($1, $2, $3)
    `)
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(email, name, phone)
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

// A guest record
type guest struct {
	Email       email
	Id          guestId
	Name        name
	PhoneNumber phoneNumber
}

// Locator for a guest record
type guestId uint8

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
	// OK +1 (555) 123-4567x890
	if !regexp.MustCompile(`^[0-9\-\(\)x ]+$`).MatchString(s) {
		return phoneNumber{}, errors.New("invalid phone number")
	}
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
