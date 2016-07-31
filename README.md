check-json-syntax - Check the syntax of a JSON files
====================================================

Nothing is more frustrating that having 6000 lines of JSON and running it in a
program to just get the message "Syntax Error" out.  No line number, no information.
That seems to be the state of JSON most of the time.  30 years ago I used 
a Fortran 66 compiler that worked that way.  It was incredibly irritating.
This program is intended to fix this for JSON.

Check the syntax of a JSON file and report decent errors.   This means output that 
shows where the syntax error is and sometimes includes suggestions on how to fix
the error.


Usage
-----

	$ check-json-syntax File1.json File2.json

Options
-------

-l Generate a listing with line numbers.

-p Pretty print the JSON if it is syntactically correct.

-D turn on debugging

Tests
-----

to run tests

	$ make test

License
------

MIT
