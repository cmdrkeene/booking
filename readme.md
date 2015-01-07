# Booking

* TODO define top level interfaces for services
* TODO test reservation service (crud for reservation records)
* TODO test billing service with dummy processor
* TODO test workflow state transitions, errors from services
* TODO basic controller/pages for workflow
* TODO workflow to offer aparment swap
* TODO admin workflow

As an admin,

- mark dates as available, unavailable
- view any offers
- accept an offer (creates an availability, booking, marks offer as accepted/hidden)
- decline an offer (deletes offer, sends decline email)

As a guest

- submit an offer
- make a booking (pay for it too)

## Domain

Workflow - a coordinator that requires payment for bookings

Guestbook
\_Guest - a record of a person staying at the hotel

Reservation - manages inventory (room(s) on dates)
\_Calendar - a master record of available/confirmed dates
\_Booking - a record of guest's payment for dates

Billing - manages bookeeping, credit cards
\_CreditCard - unsaved raw card information
\_Ledger - credits and debits for a guest's account
\_Processor - entity that manages credit cards
\_PaymentToken - persisted credit card proxy
