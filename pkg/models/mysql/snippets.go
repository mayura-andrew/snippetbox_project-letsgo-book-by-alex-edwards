package mysql

import (
	"database/sql"

	"mayuraandrew.tech/snippetbox/pkg/models"
)

//define a snippetModel type which wraps a sql.DB connection pool.

type SnippetModel struct {
	DB *sql.DB
}

// this will insert a new snippet into the database.

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// use the LastInesertId() method on the result object to get he ID of our newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// the ID returned has the type int64, so we convert ito to an int type before returning.

	return int(id), nil 
}

// this will return a spec1ific snippet based on its id.

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// use the Queryow() method on the connection pool to execute our SQL statement, passing in the untrusted id variable as the value for 
	//the placeholder parameter. this returns a pointer to a sql.Row object which 
	// holds the result fro, the database.

	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and retu
	// our own models.ErrNoRecord error instead of a Snippet object.

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// if everything went OK then return the Snippet obeject.

	return s, nil
}

// This will return the 10 most recently created snippets.

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result o
	// our query.

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.

	defer rows.Close()

	// initialize an empty slice to hold the models.Snippets objects.

	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.

	for rows.Next(){
		// create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// append it to the slice of snippets.

		snippets = append(snippets, s)
	}

		// When the rows.Next() loop has finished we call rows.Err() to retrieve any
		// error that was encountered during the iteration. It's important to
		// call this - don't assume that a successful iteration was completed
		// over the whole resultset.

	if err = rows.Err(); err != nil {
			return nil, err
	}

		// if everything went OK then return the Snippets slice

	return snippets, nil

}



