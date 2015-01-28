package booking

var formHtml = `
<html>
  <body>
    <h1>Apartment</h1>
    <h3>Book your stay</h3>
    <form action="/" method="post">
      {{ if .Errors }}
        <ul>
        {{ range .Errors }}
          <li>{{capitalize .Error}}</li>
        {{ end }}
        </ul>
      {{ end }}
      <fieldset>
        <legend>Availability</legend>
        {{range .AvailableDates}}
          {{pretty .}}
          <br />
        {{end}}
        </ul>
      </fieldset>
      <fieldset>
        <legend>Dates</legend>
        <table>
          <tr>
            <th align="left">Check In</th>
            <td><input type="text" name="checkin" placeholder="mm/dd/yyyy"/></td>
          </tr>
          <tr>
            <th align="left">Check Out</th>
            <td><input type="text" name="checkout" placeholder="mm/dd/yyyy"/></td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend>Rate</legend>
        {{range .Rates}}
        <div>
          <input name="rate" type="radio" value="{{.Name}}" />
          <b>{{.Amount}}</b>
          {{.Name}}
        </div>
        {{end}}
      </fieldset>
      <fieldset>
        <legend>Contact</legend>
        <table>
          <tr>
            <th>Name</th>
            <td><input type="text" name="name" /></td>
          </tr>
          <tr>
            <th>Email</th>
            <td><input type="text" name="email" /></td>
          </tr>
          <tr>
            <th>Phone</th>
            <td><input type="text" name="phone" /></td>
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
              <input type="text" name="card_number" />
            </td>
          </tr>
          <tr>
            <th>Month</th>
            <th>Year</th>
            <th>CVC</th>
          </tr>
          <tr>
            <td>
              <input type="text" name="card_month" size="4" />
            </td>
            <td>
              <input type="text" name="card_year" size="4"/>
            </td>
            <td>
              <input type="password" name="card_cvc" size="4"/>
            </td>
          </tr>
        </table>
      </fieldset>
      <input type="submit" value="Book" />
    </form>
  </body>
</html>`
