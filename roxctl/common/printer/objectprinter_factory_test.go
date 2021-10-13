package printer

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObjectPrinterFactory(t *testing.T) {
	cases := map[string]struct {
		defaultFormat  string
		shouldFail     bool
		errMsg         string
		printerFactory []CustomPrinterFactory
	}{
		"should fail when no CustomPrinterFactory is added": {
			defaultFormat:  "table",
			shouldFail:     true,
			errMsg:         `no custom printer factory added. You must specify at least one custom printer factory that supports the "table" output format`,
			printerFactory: []CustomPrinterFactory{nil},
		},
		"should not fail if format is supported by registered CustomPrinterFactory": {
			defaultFormat:  "table",
			shouldFail:     false,
			errMsg:         "",
			printerFactory: []CustomPrinterFactory{NewTabularPrinterFactory(false, nil, nil, false, false)},
		},
		"should fail if default output format is not supported by registered CustomPrinterFactory": {
			defaultFormat:  "table",
			shouldFail:     true,
			errMsg:         `unsupported output format used: "table". Please choose one of the supported formats: json`,
			printerFactory: []CustomPrinterFactory{NewJSONPrinterFactory(false)},
		},
		"should fail if duplicate CustomPrinterFactory is being registered": {
			defaultFormat:  "json",
			shouldFail:     true,
			errMsg:         `tried to register two printer factories which support the same output formats "json": *printer.JSONPrinterFactory and *printer.JSONPrinterFactory`,
			printerFactory: []CustomPrinterFactory{NewJSONPrinterFactory(false), NewJSONPrinterFactory(false)},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := NewObjectPrinterFactory(c.defaultFormat, c.printerFactory...)
			if c.shouldFail {
				require.Error(t, err)
				assert.Equal(t, c.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestObjectPrinterFactory_AddFlags(t *testing.T) {
	o := ObjectPrinterFactory{
		OutputFormat: "table",
		RegisteredPrinterFactories: map[string]CustomPrinterFactory{
			"json":      NewJSONPrinterFactory(false),
			"table,csv": NewTabularPrinterFactory(false, nil, nil, false, false),
		},
	}
	cmd := &cobra.Command{
		Use: "test",
	}
	o.AddFlags(cmd)
	formatFlag := cmd.Flag("format")
	assert.NotNil(t, formatFlag)
	assert.Equal(t, "table", formatFlag.DefValue)
	assert.True(t, strings.Contains(formatFlag.Usage, "json"))
	assert.True(t, strings.Contains(formatFlag.Usage, "table"))
	assert.True(t, strings.Contains(formatFlag.Usage, "csv"))
}

func TestObjectPrinterFactory_validateOutputFormat(t *testing.T) {
	cases := map[string]struct {
		o          ObjectPrinterFactory
		shouldFail bool
	}{
		"should not return an error when output format is supported": {
			o: ObjectPrinterFactory{
				OutputFormat: "table",
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"table,csv": NewTabularPrinterFactory(false, nil, nil, false, false),
					"json":      NewJSONPrinterFactory(false),
				},
			},
			shouldFail: false,
		},
		"should return an error when output format is not supported": {
			o: ObjectPrinterFactory{
				OutputFormat: "junit",
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"table,csv": NewTabularPrinterFactory(false, nil, nil, false, false),
					"json":      NewJSONPrinterFactory(false),
				},
			},
			shouldFail: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.o.validateOutputFormat()
			if c.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestObjectPrinterFactory_CreatePrinter(t *testing.T) {
	cases := map[string]struct {
		shouldFail  bool
		shouldPanic bool
		o           ObjectPrinterFactory
	}{
		"should panic when using a printer that is not yet implemented but not return an error": {
			shouldFail:  false,
			shouldPanic: true,
			o: ObjectPrinterFactory{
				OutputFormat: "table",
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"table,csv": NewTabularPrinterFactory(false, nil, nil, false, false),
				},
			},
		},
		"should return an error when the output format is not supported": {
			shouldFail:  true,
			shouldPanic: false,
			o: ObjectPrinterFactory{
				OutputFormat: "table",
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"json": NewJSONPrinterFactory(false),
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.shouldPanic {
				require.Panics(t, func() {
					_, _ = c.o.CreatePrinter()
				})
			} else {
				printer, err := c.o.CreatePrinter()
				assert.Error(t, err)
				assert.Nil(t, printer)
			}
		})
	}
}

func TestObjectPrinterFactory_validate(t *testing.T) {
	cases := map[string]struct {
		o          ObjectPrinterFactory
		shouldFail bool
		errMsg     string
	}{
		"should not fail with valid CustomPrinterFactory and valid output format": {
			o: ObjectPrinterFactory{
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"json": NewJSONPrinterFactory(false),
				},
				OutputFormat: "json",
			},
			shouldFail: false,
			errMsg:     "",
		},
		"should fail with invalid CustomPrinterFactory": {
			o: ObjectPrinterFactory{
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"table": NewTabularPrinterFactory(false, []string{"a", "b"}, []string{"a"}, false, false),
				},
				OutputFormat: "table",
			},
			shouldFail: true,
			errMsg:     "Different number of columns and JSON Path expressions specified. Make sure you specify the same number of arguments for both",
		},
		"should fail with unsupported OutputFormat": {
			o: ObjectPrinterFactory{
				RegisteredPrinterFactories: map[string]CustomPrinterFactory{
					"json": NewJSONPrinterFactory(false),
				},
				OutputFormat: "table",
			},
			shouldFail: true,
			errMsg:     `unsupported output format used: "table". Please choose one of the supported formats: json`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.o.validate()
			if c.shouldFail {
				require.Error(t, err)
				assert.Equal(t, c.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}