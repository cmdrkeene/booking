package booking

// func xTestHandler(t *testing.T) {
// 	db := testDB()
// 	defer db.Close()
// 	var calendar Calendar
// 	var handler Handler
// 	err := inject.Populate(
// 		db,
// 		&calendar,
// 		&handler,
// 	)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	calendar.Add(
// 		date.New(2015, 1, 1),
// 		date.New(2015, 1, 2),
// 		date.New(2015, 1, 3),
// 	)

// 	// get lists available dates
// 	w := httptest.NewRecorder()
// 	r, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	handler.ServeHTTP(w, r)

// 	if w.Code != http.StatusOK {
// 		t.Error("want", http.StatusOK)
// 		t.Error("got ", w.Code)
// 	}

// 	body := w.Body.Bytes()
// 	dates := []string{
// 		"January 1, 2015",
// 		"January 2, 2015",
// 		"January 3, 2015",
// 	}
// 	for _, s := range dates {
// 		if !bytes.Contains(body, []byte(s)) {
// 			t.Error("want", s)
// 			t.Error("got ", string(body))
// 		}
// 	}

// 	// post registers, pays, and books dates
// 	vals := url.Values{}
// 	vals.Add(formKeyCardCVC, "123")
// 	vals.Add(formKeyCardMonth, "01")
// 	vals.Add(formKeyCardNumber, "1111222233334444")
// 	vals.Add(formKeyCardYear, "15")
// 	vals.Add(formKeyCheckin, "2015-01-01")
// 	vals.Add(formKeyCheckout, "2015-01-02")
// 	vals.Add(formKeyCheckout, "2015-01-03")
// 	vals.Add(formKeyEmail, "a@b")
// 	vals.Add(formKeyName, "a b")

// 	w = httptest.NewRecorder()
// 	r, err = http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	r.Header.Set(
// 		"Content-Type",
// 		"application/x-www-form-urlencoded; param=value",
// 	)

// 	handler.ServeHTTP(w, r)

// 	if w.Code != 201 {
// 		t.Error("want", 201)
// 		t.Error("got ", w.Code)
// 	}
// }
