---
title: "CSV Boilerplate for Go"
date: 2020-09-27
---

**tl;dr:** Here is [my boilerplate](https://github.com/felixge/dump/blob/master/csv-boilerplate/main.go) for reading and writing CSV files in Go. I hope you'll find it useful.

Every once in a while, I have to read and/or write CSV files in Go. In theory that's easy. One just has to include [encoding/csv](https://golang.org/pkg/encoding/csv/) from the standard library, and write a bit of boilerplate to map between a Go struct and the CSV records.

However, the resulting code has often left me unsatisifed for a variety of reasons:

1. I spend too much time reinventing the wheel, playing around with different ways to structure the code.
2. The code is annoying to maintain, e.g. changing the output column order often requires modifications in at least 4 places: the struct definition, header writing, record writing, and record reading.
3. The result is not very robust when it comes to bad CSV input and doesn't handle problems such as missing, unexpected, or misordered columns or headers.
4. Error messages often end up missing important information such as row number, column name, etc..
5. How do I write good tests for this that will be easy to maintain in the future?

One way to deal with this is to use a high level library, and you should certainly consider that. But the reality is often an ugly mess. You might find that the library of your choice doesn't handle a particular edge case, e.g. more than one header row being present (my bank loves putting in 6 header lines). Or the error messages you're getting are not helpful enough for your users. Or the error handling is too strict. Or the reflection based mapping is too magic. Or the performance sucks. I could go on and on, but you get the point - dealing with CSV is nasty and calling it a format is a bit of an insult to the idea of a format. From my point of view, writing a library for this problem is simply not possible without creating a horribly overengineered mess or ignoring many edge cases.

So in the end it seems like you have to chose between two suboptimal options. You can try out a bunch of libraries and use and/or fork the one closest to your needs. Or you can roll your own boilerplate for the millionth time, only to look back at it with horrible disgust shortly after.

Your choice should probably depend on the craziness of your input, and I won't be able to make it for you. But I've decided that I want to make the boilerplate route a bit less annoying for myself in the future. To accomplish this I've created a [boilerplate template](https://github.com/felixge/dump/blob/master/csv-boilerplate/main.go) that I'll copy and paste for my future needs, adjusting it as needed. For the rest of the blog post I want to walk you through my thought process in designing it.

To begin, I decided to use real world CSV as my example, containing bits of the typical ugliness one has to deal with. I quickly [found](https://perso.telecom-paristech.fr/eagan/class/igr204/datasets) a file called [cars.csv](https://github.com/felixge/dump/blob/master/csv-boilerplate/cars.csv) which looks like this:

```csv
Car;MPG;Cylinders;Displacement;Horsepower;Weight;Acceleration;Model;Origin
STRING;DOUBLE;INT;DOUBLE;DOUBLE;DOUBLE;DOUBLE;INT;CAT
Chevrolet Chevelle Malibu;18.0;8;307.0;130.0;3504.;12.0;70;US
Buick Skylark 320;15.0;8;350.0;165.0;3693.;11.5;70;US
Plymouth Satellite;18.0;8;318.0;150.0;3436.;11.0;70;US
...
```

As you can see, there are already two ugly bits to be examined here: semicolons instead of commas as a separator, and a second header line that seems to contain type information for whatever reason. Not great, not terrible, but pretty typical for what's out there.

The goals for my boilerplate is to read and write the file format in a way that solves the three problems outlined in the beginning. It should do so without an excessive amount of abstraction, and certainly no reflection, yet as elegantly as possible. There should also be an intermediate Go struct to hold the data in memory as shown below.

```go
type Car struct {
	Car          string
	MPG          float64
	Cylinders    int64
	Displacement float64
	Horsepower   float64
	Weight       float64
	Acceleration float64
	Model        int64
	Origin       string
}
```

My first naive attempt is usually to write some code that defines the headers and struct mapping in a way that is separate from the CSV encoding/decoding itself. An example of that can be seen below:

```go
var carHeaders = []string{
  "Car",
  "MPG",
  "Cylinders",
  // remaining columns ...
}

func (c *Car) UnmarshalRecord(record []string) (err error) {
	if got, want := len(record), len(carHeaders); got != want {
		return fmt.Errorf("bad column number: got=%d want=%d", got, want)
	}
	c.Car = record[0]
	if c.MPG, err = strconv.ParseFloat(record[1], 64); err != nil {
		return fmt.Errorf("column=%q: %w", carHeaders[1], err)
	} else if c.Cylinders, err = strconv.ParseInt(record[2], 10, 64); err != nil {
		return fmt.Errorf("column=%q: %w", carHeaders[2], err)
	}
	// remaining columns ...
	return nil
}

func (c *Car) MarshalRecord() ([]string, error) {
	return []string{
		c.Car,
		fmt.Sprintf("%f", c.MPG),
		fmt.Sprintf("%d", c.Cylinders),
		// remaining columns ...
	}, nil
}
```

This works, but adding, removing, or reordering any of the columns requires us to modify our code in 4 to 5 different places, which we wanted to avoid. So let's try some abstraction that allows us to declaratively specify our columns and how to marshal/unmarshal them:

```go
type carColumn struct {
	Name           string
	Type           string
	UnmarshalValue func(*Car, string) error
	MarshalValue   func(*Car) (string, error)
}

var carColumns = []carColumn{
	{
		"Car",
		"STRING",
		func(c *Car, val string) error {
			c.Car = val
			return nil
		},
		func(c *Car) (string, error) {
			return c.Car, nil
		},
	},
	{
		"MPG",
		"DOUBLE",
		func(c *Car, val string) (err error) {
			c.MPG, err = strconv.ParseFloat(val, 64)
			return
		},
		func(c *Car) (string, error) {
			return fmt.Sprintf("%f", c.MPG), nil
		},
	},
	{
		"Cylinders",
		"INT",
		func(c *Car, val string) (err error) {
			c.Cylinders, err = strconv.ParseInt(val, 10, 64)
			return
		},
		func(c *Car) (string, error) {
			return fmt.Sprintf("%d", c.Cylinders), nil
		},
	},
	// remaining columns ...
}
```

This is of course a bit verbose, but on the other hand it solves our problem. Reordering our columns is now a single cut & paste operation, and adding or removing a column has gotten a lot easier as well. The approach even gives us a convenient way to put the `Type` information found in the second line of the CSV file, so we can easily duplicate this quirk when writing our own CSV files later on.

Of course we still need to update our `UnmarshalRecord` and `MarshalRecord` from earlier. The good news is that we're unlikely to ever have to modify this code again:

```go
func (c *Car) UnmarshalRecord(record []string) error {
	if got, want := len(record), len(carColumns); got != want {
		return fmt.Errorf("bad number of columns: got=%d want=%d", got, want)
	}
	for i, col := range carColumns {
		if err := col.UnmarshalValue(c, record[i]); err != nil {
			return fmt.Errorf("column=%q: %w", col.Name, err)
		}
	}
	return nil
}

func (c *Car) MarshalRecord() ([]string, error) {
	record := make([]string, len(carColumns))
	for i, col := range carColumns {
		val, err := col.MarshalValue(c)
		if err != nil {
			return nil, err
		}
		record[i] = val
	}
	return record, nil
}
```

Now it's time to implement the CSV decoding itself. My boilerplate implements this via the `ReadCarsCSV` function that can be seen below.

```go
const carComma = ';'

func ReadCarsCSV(r io.Reader) ([]*Car, error) {
	cr := csv.NewReader(r)
	cr.Comma = carComma

	records, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	var cars []*Car
	for i, record := range records {
		if got, want := len(record), len(carColumns); got != want {
			return nil, fmt.Errorf("row=%d: bad number of columns: got=%d want=%d", i+1, got, want)
		}

		switch i {
		case 0:
			for i, got := range record {
				if want := carColumns[i].Name; got != want {
					return nil, fmt.Errorf("unexpected header column %d: got=%q want=%q", i, got, want)
				}
			}
		case 1:
			for i, got := range record {
				if want := carColumns[i].Type; got != want {
					return nil, fmt.Errorf("unexpected type column %d: got=%q want=%q", i, got, want)
				}
			}
		default:
			car := &Car{}
			if err := car.UnmarshalRecord(record); err != nil {
				return nil, fmt.Errorf("row=%d: %w", i+1, err)
			}
			cars = append(cars, car)
		}
	}

	return cars, nil
}
```

The code above accomplishes the main task of returning all cars read from the given `io.Reader`, but it also validates the number of columns, as well as the header and type names found in the first two rows. Error messages should also be good, including both the row number as well the the offending column name in case something goes wrong for one of the records.

It should also be easy to modify. For example if you want to ignore the second header line during reading, just remove the code. Or if instead of requiring all headers to have a fixed position, you could dynamically discover their position from the input by matching their names against the names in the `carColumns` slice and then reorder the elements in the record accordingly.

Or you might decide you need a streamining interface to lower memory usage and GC pressure like this:

```go
type CarReader struct {
	r io.Reader
	// ...
}

func (r *CarReader) Read(c *Car) error {
	// ...
}
```

The possibilities are endless, and you won't find yourself backed into a corner by the choice of your library.

Next up is converting our data back to CSV. My solution to that looks like this:


```go
func WriteCarsCSV(w io.Writer, cars []*Car) error {
	cw := csv.NewWriter(w)
	cw.Comma = carComma

	header := make([]string, len(carColumns))
	types := make([]string, len(carColumns))
	for i, col := range carColumns {
		header[i] = col.Name
		types[i] = col.Type
	}
	cw.Write(header)
	cw.Write(types)

	for _, car := range cars {
		record, err := car.MarshalRecord()
		if err != nil {
			return err
		}
		cw.Write(record)
	}
	cw.Flush()
	return cw.Error()
}
```

As you can see, we take care of writing out the quirky second header line. We also save some code by not checking the errors for every `cw.Write()` call and handle them by returning `cw.Error()` at the end instead.

Now we should think about testing. I'm a big fan of the [80/20 rule](https://en.wikipedia.org/wiki/Pareto_principle) when it comes to testing, so below is a test that covers the happy path very efficiently.

```go
func TestCSVReadWriteCycle(t *testing.T) {
	in := strings.TrimSpace(`
Car;MPG;Cylinders;Displacement;Horsepower;Weight;Acceleration;Model;Origin
STRING;DOUBLE;INT;DOUBLE;DOUBLE;DOUBLE;DOUBLE;INT;CAT
Chevrolet Chevelle Malibu;18.0;8;307.0;130.0;3504.;12.0;70;US
Buick Skylark 320;15.0;8;350.0;165.0;3693.;11.5;70;US
Plymouth Satellite;18.0;8;318.0;150.0;3436.;11.0;70;US
`)
	wantOut := strings.TrimSpace(`
Car;MPG;Cylinders;Displacement;Horsepower;Weight;Acceleration;Model;Origin
STRING;DOUBLE;INT;DOUBLE;DOUBLE;DOUBLE;DOUBLE;INT;CAT
Chevrolet Chevelle Malibu;18.000000;8;307.000000;130.000000;3504.000000;12.000000;70;US
Buick Skylark 320;15.000000;8;350.000000;165.000000;3693.000000;11.500000;70;US
Plymouth Satellite;18.000000;8;318.000000;150.000000;3436.000000;11.000000;70;US
`)

	for i := 0; i < 2; i++ {
		cars, err := ReadCarsCSV(strings.NewReader(in))
		if err != nil {
			t.Fatal(err)
		}

		buf := &bytes.Buffer{}
		if err := WriteCarsCSV(buf, cars); err != nil {
			t.Fatal(err)
		}
		gotOut := strings.TrimSpace(buf.String())
		if gotOut != wantOut {
			t.Fatalf("\ngot:\n%s\nwant:\n%s", gotOut, wantOut)
		}
		in = gotOut
	}
}
```

The test first reads the provided `in` CSV, and makes sure that converting it back to CSV results in `wantOut` which is a little different because of the way Go formats our floating point numbers. The test then proceeds to treat the output from the first iteration as the input for the second one. This makes sure that our implementation can read the original format, as well as the slightly different output it produces itself.

Of course we could also aim to reproduce the float formatting quirks from the original file. At first glance it seems like all floats are formatted with one digit after the period. However, a closer look reveals that the `Weight` column has a period, but skips the following digit, e.g. `3504.`. I've decided to not go down this rabbit hole, but it should be clear that our boilerplate approach puts us in a great position for dealing with CSV quirks like this.

Depending on how serious you are, you probably also want to test a few error cases and maybe even throw some fuzzing at this. But I've decided the test above is good enough for my boilerplate, so you'll have to do this part yourself.

Anyway, thanks for taking the time to read this post. I make no claim that my approach for dealing with CSV is superior to all alternatives, but I think you could probably do a lot worse than starting with [my boilerplate](https://github.com/felixge/dump/blob/master/csv-boilerplate/main.go).

That being said, I'd love to hear your thoughts, especially if you have ideas for improving it further. Or let me know if you see good alternatives to the idea of copy & pasting a bunch of boilerplate. It's not like I love it, but as far as I can tell it hits a sweet spot for working with Go in this case.
