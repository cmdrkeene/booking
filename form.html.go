package booking

import (
	"html/template"
	"unicode"

	"github.com/cmdrkeene/booking/pkg/date"
)

var templateForm *template.Template

var formHelpers = template.FuncMap{
	"capitalize": func(s string) string {
		a := []rune(s)
		a[0] = unicode.ToUpper(a[0])
		return string(a)
	},
	"pretty": func(d date.Date) string {
		return d.Format(date.Pretty)
	},
	"equals": func(a, b string) bool {
		return a == b
	},
}

func init() {
	templateForm = template.Must(
		template.New("form.html.go").Funcs(formHelpers).Parse(templateFormSrc),
	)
}

const templateFormSrc = `
<html>
  <head>
    <style type="text/css">
      ul.error li {
        color: red;
      }
      input.error {
        border: 1px solid red;
      }
    </style>
  </head>
  <body>
    <h1>Apartment</h1>
    <h3>Book your stay</h3>
    <form action="/" method="post">
      <!-- All Errors -->
      {{if .Errors}}
        <ul class="error">
        {{range $field, $error := .Errors}}
          <li>{{$field}} {{$error}}</li>
        {{end}}
        </ul>
      {{end}}

      <!-- Available Dates -->
      <fieldset>
        <legend>Available Dates</legend>
        
        {{range .AvailableDates}}
          {{pretty .}}
          <br />
        {{end}}
        </ul>
      </fieldset>
      
      <!-- Your Dates -->
      <fieldset>
        <legend>Your Dates</legend>

        <label>Check in</label>
        <input 
          type="text" 
          name="Checkin" 
          placeholder="mm/dd/yyyy" 
          value="{{.Checkin}}" 
          {{with .Errors.Checkin}}
            class="error"
          {{end}}
        />

        <label>Check out</label>
        <input 
          type="text" 
          name="Checkout" 
          placeholder="mm/dd/yyyy" 
          value="{{.Checkout}}"
          {{with .Errors.Checkout}}
            class="error"
          {{end}}
        />
      </fieldset>

      <!-- Rate -->
      <fieldset>
        <legend>Rate</legend>
        
        <!-- errors -->
        {{with .Errors.Rates}}
          {{.}}
        {{end}}

        <!-- all rates  -->
        {{$currentRate := .Rate}}
        {{range .Rates}}
        <div>
          <input
            name="Rate" 
            type="radio" 
            value="{{.Name}}" 
            {{if equals $currentRate .Name}}
              checked
            {{end}}
            />
          <b>{{.Amount}}</b>
          {{.Name}}
        </div>
        {{end}}
      </fieldset>
      <fieldset>
        <legend>Guest</legend>
        <table>
          <tr>
            <th>Name</th>
            <td>
              <input 
                type="text" 
                name="Name" 
                value="{{.Name}}"
                {{with .Errors.Name}}
                  class="error"
                {{end}}
              />
            </td>
          </tr>
          <tr>
            <th>Email</th>
            <td>
              <input
                type="text"
                name="Email"
                value="{{.Email}}"
                {{with .Errors.Email}}
                  class="error"
                {{end}}
              />
            </td>
          </tr>
          <tr>
            <th>Phone</th>
            <td>
              <input
                type="text"
                name="Phone"
                value="{{.Phone}}"
                {{with .Errors.Phone}}
                  class="error"
                {{end}}
              />
            </td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend>Credit Card</legend>
        <table>
          <tr>
            <th colspan="3" align="left">Number</th>
          </tr>
          <tr>
            <td colspan="3">
              <input 
                type="text" 
                name="CardNumber" 
                value="{{.CardNumber}}"
                {{with .Errors.CardNumber}}
                  class="error"
                {{end}}
              />
            </td>
          </tr>
          <tr>
            <th>Month</th>
            <th>Year</th>
            <th>CVC</th>
          </tr>
          <tr>
            <td>
              <input 
                type="text" 
                name="CardMonth" 
                size="4" 
                value="{{.CardMonth}}"
                {{with .Errors.CardMonth}}
                  class="error"
                {{end}}
              />
            </td>
            <td>
              <input
                type="text"
                name="CardYear"
                size="4"
                value="{{.CardYear}}"
                {{with .Errors.CardYear}}
                  class="error"
                {{end}}
              />
            </td>
            <td>
              <input
                type="password" 
                name="CardCVC" 
                size="4" 
                value="{{.CardCVC}}"
                {{with .Errors.CardCVC}}
                  class="error"
                {{end}}
              />
            </td>
          </tr>
        </table>
      </fieldset>
      <input type="submit" value="Book" />
    </form>
  </body>
</html>
`
